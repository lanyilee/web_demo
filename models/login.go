package models

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/astaxie/beego/context"
)

type UserInfo struct {
	UserID       string
	UserName     string
	Role         string
	HarborToken  string
	RancherToken string
	DroneToken   string
	BaseToken    string
}

type LoginSpec struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthInfterface interface {
	CreateToken(info UserInfo) string
	//CreateDroneToken(info UserInfo)string
	GetUserInfo(tokenStr string) (UserInfo, error)
	ParseFromRequestToken(req *http.Request) (UserInfo, error)
	JwtAuthFilter(ctx *context.Context)
	JwtAuthFilterGin(ctx *gin.Context)
	ParseFromSecret(req *http.Request) (UserInfo, error)
}

type OIDCInterface interface {
	Callback() http.Handler
	Login() http.Handler
	HandleLogin(w http.ResponseWriter, r *http.Request)
}
