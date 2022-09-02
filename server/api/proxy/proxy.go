package proxy

import (
	"crypto/tls"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"webase-server/models"
)

type proxy struct {
	auth   models.AuthInfterface
	addr   string
	prefix string
}

func (p *proxy) proxy(req *http.Request) error {
	_, err := p.auth.ParseFromRequestToken(req)
	if err != nil {
		return err
	}
	path := req.URL.String()
	index := strings.Index(path, p.prefix)
	destPath := path[index+len(p.prefix):]
	destURL, _ := url.Parse(p.addr + destPath)
	remote, err := url.Parse(p.addr)
	if err != nil {
		zap.S().Info(err)
		return err
	}
	req.Host = remote.Host
	req.URL = destURL
	return nil
}

func (p *proxy) proxyNoAuth(req *http.Request) error {
	path := req.URL.String()
	index := strings.Index(path, p.prefix)
	destPath := path[index+len(p.prefix):]
	destURL, _ := url.Parse(p.addr + destPath)
	remote, err := url.Parse(p.addr)
	if err != nil {
		zap.S().Info(err)
		return err
	}
	req.Host = remote.Host
	req.URL = destURL
	return nil
}

func New(auth models.AuthInfterface, addr, prefix string) http.Handler {
	p := &proxy{auth, addr, prefix}
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if err := p.proxy(req); err != nil {
				zap.S().Infof("Failed to proxy: %v", err)
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func NewNoAuth(auth models.AuthInfterface, addr, prefix string) http.Handler {
	p := &proxy{auth, addr, prefix}
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if err := p.proxyNoAuth(req); err != nil {
				zap.S().Infof("Failed to proxy: %v", err)
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
