package v3

import (
	"encoding/base64"
	"go.uber.org/zap"
	"webase-server/models"
	"webase-server/server/api/base"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	base.BaseV3Controller
}

type userToken struct {
	IsBindDrone bool   `json:"is_bind_drone"`
	Token       string `json:"token"`
	models.LocalUser
}

func (l *LoginController) Login(c *gin.Context) {
	var p models.LoginSpec
	p = models.LoginSpec{}
	//p.Username = c.PostForm("username")
	//p.Password = c.PostForm("password")
	err := c.ShouldBindJSON(&p)
	if err != nil {
		l.Logger.Error(err)
		l.SetErrMsg(c, err, 500)
		return
	}

	u, err := l.APIManager.Store.User.Login(p.Username, p.Password)
	if err != nil {
		l.Logger.Error(err)
		l.SetErrMsg(c, err, 500)
		return
	}
	userOut := userToken{
		LocalUser: *u,
	}
	baseToken := base64.StdEncoding.EncodeToString([]byte(p.Username + ":" + p.Password))
	info := models.UserInfo{
		UserID:      u.ID,
		UserName:    u.Username,
		Role:        u.Role,
		HarborToken: baseToken,
		BaseToken:   baseToken,
	}
	//if u.Role == models.AdminRole {
	//	info.HarborToken = base64.StdEncoding.EncodeToString([]byte(l.APIManager.Config.Harbor.Username + ":" + l.APIManager.Config.Harbor.Password))
	//}
	binds, err := l.APIManager.Store.UserBind.FindBindsByUser(u.ID)
	if err != nil {
		l.Logger.Error(err)
		l.SetErrMsg(c, err, 500)
		return
	}
	//add
	println(&binds)

	userOut.Token = l.APIManager.Auth.CreateToken(info)
	//c.SetCookie(models.CookiePath, userOut.Token, 28800, "/", "", true, false)
	l.SetResult(c, userOut, 200)
}

func (l *LoginController) Logout(c *gin.Context) {

	path := "xxx"
	//path := l.APIManager.Config.OIDC.IssuerURL + "/end?redirect_uri=" + l.APIManager.Config.Host + "/ui"
	zap.S().Debug(path)
	c.Redirect(302, path)
}
