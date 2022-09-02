package routers

import (
	"github.com/gin-gonic/gin"
	"webase-server/models"
	"webase-server/server/api/base"
	v3 "webase-server/server/api/v3"
)

// NewV3 初始化 v3 接口路由
func NewV3(router *gin.Engine, apimgr *models.APIManager) {
	basec := base.BaseV3Controller{
		Logger:     apimgr.Logger,
		APIManager: apimgr,
	}
	loginCtrl := &v3.LoginController{BaseV3Controller: basec}
	userCtrl := &v3.UserController{BaseV3Controller: basec}



	publicv3 := router.Group("/webase/public/v2")
	{
		publicv3.POST("/login", loginCtrl.Login)
		publicv3.GET("/logout", loginCtrl.Logout)
	}

	apiv3 := router.Group("/webase/api/v2/")
	apiv3.Use(apimgr.Auth.JwtAuthFilterGin)
	{
		userGroup := apiv3.Group("/user")
		userGroup.POST("", userCtrl.Create)
		userGroup.GET("", userCtrl.Current)
		userGroup.PUT("", userCtrl.UpdateCurrent)

		usersGroup := apiv3.Group("/users")
		usersGroup.GET("", userCtrl.List)
		usersGroup.GET("/:uid", userCtrl.Get)
		usersGroup.DELETE("/:uid", userCtrl.Delete)
		usersGroup.PUT("/:uid", userCtrl.Update)


	}




	//app metrics
	//router.Any("/app/metrics", gin.WrapH(promhttp.Handler()))

}
