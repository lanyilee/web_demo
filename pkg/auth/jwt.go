package auth

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"webase-server/models"
	"webase-server/pkg/utils/crypto"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type jwtClient struct {
	secretKey string
	expire    time.Duration
	store     *models.Store
}

func New(secretKey string, expire time.Duration, store *models.Store) models.AuthInfterface {
	if secretKey == "" {
		secretKey = "ktFEzmcWRr91pQPZ"
	}
	return &jwtClient{secretKey, expire, store}
}

type serverClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func (c *jwtClient) encode(info models.UserInfo) string {
	data, _ := json.Marshal(info)
	return hex.EncodeToString(crypto.AESEncrypt(data, []byte(c.secretKey)))
}

func (c *jwtClient) decode(s string) (models.UserInfo, error) {
	var info models.UserInfo
	data, err := hex.DecodeString(s)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(crypto.AESDecrypt(data, []byte(c.secretKey)), &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

func (c *jwtClient) CreateToken(info models.UserInfo) string {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := serverClaims{
		User: c.encode(info),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "webase",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(c.secretKey))
	return signedToken
}

func (c *jwtClient) GetUserInfo(tokenStr string) (models.UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &serverClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(c.secretKey), nil
	})
	if err != nil {
		return models.UserInfo{}, err
	}

	if claims, ok := token.Claims.(*serverClaims); ok && token.Valid {
		info, err := c.decode(claims.User)
		if err != nil {
			return models.UserInfo{}, err
		}
		return info, nil
	}

	return models.UserInfo{}, errors.New("must auth")

}

func (c *jwtClient) ParseFromRequestToken(req *http.Request) (models.UserInfo, error) {
	tokenStr := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if tokenStr == "" {
		tokenStr = req.FormValue("token")
	}
	info, err := c.GetUserInfo(tokenStr)
	return info, err
}

func (c *jwtClient) JwtAuthFilter(ctx *context.Context) {
	if ctx.Request.FormValue("token") != "" {
		info, err := c.GetUserInfo(ctx.Request.FormValue("token"))
		if err == nil {
			ctx.Input.SetData("user_info", info)
			return
		}
	}
	token, err := request.ParseFromRequestWithClaims(ctx.Request,
		request.AuthorizationHeaderExtractor,
		&serverClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(c.secretKey), nil
		})
	if err != nil {
		ctx.Output.Status = 401
		ctx.Output.JSON(err, false, false)
		return
	}
	info, err := c.decode(token.Claims.(*serverClaims).User)
	if err != nil {
		ctx.Output.Status = 401
		zap.S().Error(err)
		ctx.Output.JSON(err, false, false)
		return
	}
	ctx.Input.SetData("user_info", info)
}

//add gin jwtauth
func (c *jwtClient) JwtAuthFilterGin(ctx *gin.Context) {
	//coo, err := ctx.Request.Cookie(models.CookiePath)
	//if err != nil {
	//	fmt.Println("cookie error")
	//}
	//fmt.Println("cookie 过期时间")
	//fmt.Println(coo.Expires.Date())

	if ctx.Request.FormValue("token") != "" {
		info, err := c.GetUserInfo(ctx.Request.FormValue("token"))
		fmt.Print(info)
		if err == nil {
			ctx.Set("user_info", info)
			return
		}
	}

	//fmt.Println("Authorization:" + ctx.Request.Header.Get("Authorization"))
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(c.secretKey), nil
		}, request.WithClaims(&serverClaims{}))
	if err != nil {
		ctx.JSON(401, err)
		return
	}
	info, err := c.decode(token.Claims.(*serverClaims).User)
	if err != nil {
		zap.S().Error(err)
		ctx.JSON(401, err)
		return
	}
	ctx.Set("user_info", info)
	ctx.Next()
}

var SecretKey = "ssfesssfgss"

