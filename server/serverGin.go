package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"webase-server/models"
	"webase-server/server/middlewaresGin"
	"webase-server/server/routers"
)

func StartGin(apimgr *models.APIManager) {
	r := gin.New()
	r.Use(middlewaresGin.Logger())
	r.Use(middlewaresGin.Cors)

	//r.Handle("GET", "/hello", func(context *gin.Context) {
	//	//panic(apimgr)
	//	time.Sleep(time.Second * 2)
	//	context.JSON(200, "hello world")
	//})
	//r.Run(":8888")

	routers.NewV3(r, apimgr)
	r.Run(":" + strconv.Itoa(apimgr.Config.Port))

}
