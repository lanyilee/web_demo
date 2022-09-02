package oidc

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"go.uber.org/zap"

	"webase-server/models"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type client struct {
	clientID       string
	clientSecret   string
	redirectURI    string
	issuerURL      string
	rootCAs        string
	verifier       *oidc.IDTokenVerifier
	provider       *oidc.Provider
	offlineAsScope bool
	client         *http.Client
	auth           models.AuthInfterface
	store          *models.Store
	authExporter   models.AuthExporterInterface
}

//New OIDC client
func New(store *models.Store, auth models.AuthInfterface, clientID, clientSecret, issuerURL, redirectURI, rootCAs string, authExporter models.AuthExporterInterface) (models.OIDCInterface, error) {
	c := &client{
		clientID:       clientID,
		clientSecret:   clientSecret,
		issuerURL:      issuerURL,
		redirectURI:    redirectURI,
		rootCAs:        rootCAs,
		offlineAsScope: true,
		auth:           auth,
		store:          store,
		authExporter:   authExporter,
	}
	err := c.init()
	return c, err
}

type callbackH struct {
	c *client
}

func (h *callbackH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.HandleCallback(w, r)
}

func (c *client) Callback() http.Handler {
	return &callbackH{c}
}

type loginH struct {
	c *client
}

func (h *loginH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.HandleLogin(w, r)
}

func (c *client) Login() http.Handler {
	return &loginH{c}
}

func (c *client) init() error {
	if c.rootCAs != "" {
		client, err := httpClientForRootCAs(c.rootCAs)
		if err != nil {
			return err
		}
		c.client = client
	}
	if c.client == nil {
		c.client = models.DefaultHTTPClient
	}
	ctx := oidc.ClientContext(context.Background(), c.client)
	provider, err := oidc.NewProvider(ctx, c.issuerURL)
	if err != nil {
		return fmt.Errorf("failed to query provider %q: %v", c.issuerURL, err)
	}
	var s struct {
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := provider.Claims(&s); err != nil {
		return fmt.Errorf("failed to parse provider scopes_supported: %v", err)
	}

	if len(s.ScopesSupported) == 0 {
		c.offlineAsScope = true
	} else {
		c.offlineAsScope = func() bool {
			for _, scope := range s.ScopesSupported {
				if scope == oidc.ScopeOfflineAccess {
					return true
				}
			}
			return false
		}()
	}

	c.provider = provider
	c.verifier = provider.Verifier(&oidc.Config{ClientID: c.clientID})
	return nil
}

// return an HTTP client which trusts the provided root CAs.
func httpClientForRootCAs(rootCAs string) (*http.Client, error) {
	tlsConfig := tls.Config{RootCAs: x509.NewCertPool()}
	rootCABytes, err := ioutil.ReadFile(rootCAs)
	if err != nil {
		return nil, fmt.Errorf("failed to read root-ca: %v", err)
	}
	if !tlsConfig.RootCAs.AppendCertsFromPEM(rootCABytes) {
		return nil, fmt.Errorf("no certs found in root CA file %q", rootCAs)
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}, nil
}

type debugTransport struct {
	t http.RoundTripper
}

func (d debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	zap.S().Infof("%s", reqDump)

	resp, err := d.t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	zap.S().Infof("%s", respDump)
	return resp, nil
}

func (c *client) oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Endpoint:     c.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  c.redirectURI,
	}
}

//需要携带scope、response_type、client_id、redirect_uri、state
func (c *client) HandleLogin(w http.ResponseWriter, r *http.Request) {
	authCodeURL := ""
	scopes := []string{"openid", "profile", "email"}
	appState := createRandomString(5)
	if c.offlineAsScope {
		scopes = append(scopes, "offline_access")
		authCodeURL = c.oauth2Config(scopes).AuthCodeURL(appState)
	} else {
		authCodeURL = c.oauth2Config(scopes).AuthCodeURL(appState, oauth2.AccessTypeOffline)
	}
	cookie := &http.Cookie{
		Name:   "app_state",
		Value:  appState,
		Path:   "/",
		MaxAge: 600,
	}
	http.SetCookie(w, cookie)
	//去idaas中认证
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (c *client)HandleCallback(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token *oauth2.Token
	)

	ctx := oidc.ClientContext(r.Context(), c.client)
	oauth2Config := c.oauth2Config(nil)
	switch r.Method {
	case http.MethodGet:
		// Authorization redirect callback from OAuth2 auth flow.
		if errMsg := r.FormValue("error"); errMsg != "" {
			http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
			c.authExporter.HandleAuthFailCounter(400)
			return
		}
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
			c.authExporter.HandleAuthFailCounter(400)
			return
		}
		cookie, err := r.Cookie("app_state")
		if err != nil {
			http.Error(w, fmt.Sprintf("登录过期，请重新登录！"), http.StatusBadRequest)
			c.authExporter.HandleAuthFailCounter(400)
			return
		}
		if state := r.FormValue("state"); state != cookie.Value {
			http.Error(w, fmt.Sprintf("expected state %q got %q", cookie.Value, state), http.StatusBadRequest)
			c.authExporter.HandleAuthFailCounter(400)
			return
		}
		token, err = oauth2Config.Exchange(ctx, code)
	case http.MethodPost:
		// Form request from frontend to refresh a token.
		refresh := r.FormValue("refresh_token")
		if refresh == "" {
			http.Error(w, fmt.Sprintf("no refresh_token in request: %q", r.Form), http.StatusBadRequest)
			c.authExporter.HandleAuthFailCounter(400)
			return
		}
		t := &oauth2.Token{
			RefreshToken: refresh,
			Expiry:       time.Now().Add(-time.Hour),
		}
		token, err = oauth2Config.TokenSource(ctx, t).Token()
	default:
		http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusBadRequest)
		c.authExporter.HandleAuthFailCounter(400)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}

	idToken, err := c.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify ID token: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	// accessToken, ok := token.Extra("access_token").(string)
	// if !ok {
	// 	http.Error(w, "no access_token in token response", http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println(accessToken, rawIDToken)
	// fmt.Println(rawIDToken, accessToken, token.RefreshToken)
	var claims json.RawMessage
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, fmt.Sprintf("error decoding ID token claims: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	u := userInfo{}
	if err := json.Unmarshal([]byte(claims), &u); err != nil {
		http.Error(w, fmt.Sprintf("error indenting ID token claims: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	user, exist, err := c.store.User.GetUserByUsername(u.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("found user info error: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	if !exist {
		http.Error(w, " not found user info", http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	fmt.Println(user.Name)
	binds, err := c.store.UserBind.FindBindsByUser(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("found user info error: %v", err), http.StatusInternalServerError)
		c.authExporter.HandleAuthFailCounter(500)
		return
	}
	info := models.UserInfo{
		Role:     user.Role,
		UserID:   user.ID,
		UserName: user.Username,
	}
	for _, bind := range binds {
		if bind.Provider == models.DroneProvider {
			info.DroneToken = bind.ProviderID
			cookie := &http.Cookie{
				Name:   models.DroneCookiePath,
				Value:  bind.ProviderID,
				Path:   "/",
				MaxAge: 28800,
			}
			http.SetCookie(w, cookie)
		}
	}
	tokenStr := c.auth.CreateToken(info)
	cookie := &http.Cookie{
		Name:   models.CookiePath,
		Value:  tokenStr,
		Path:   "/",
		MaxAge: 28800,
	}
	http.SetCookie(w, cookie)
	fmt.Println("handler: " + tokenStr)
	c.authExporter.HandleAuthSuccessCounter(200)
	http.Redirect(w, r, "/ui", 303)

}

type userInfo struct {
	Username string `json:"preferred_username"`
	Name     string `json:"name"`
}

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
