package netutils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

//获取gin requestbody 并放回去
func GetGinRequestBody(ctx *gin.Context)[]byte{
	var bodyBytes []byte // 我们需要的body内容
	// 从原有Request.Body读取
	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {

	}
	// 新建缓冲区并替换原有Request.body
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 当前函数可以使用body内容
	return bodyBytes
}

//获取url路径组件信息
func  GetComponent(ctx *gin.Context) string {
	//fmt.Println(ctx.Request.URL.Path)
	path := ctx.Request.URL.Path
	us := strings.Split(path, "/")
	if strings.HasPrefix(path,"/webase/api/v2/") && len(us)>4{
		if us[4]=="app"{
			return "helm"
		}else if us[4]=="application"{
			return "docker-compose"
		}
	}
	return us[1]
}

//分析url路径，并返回是否是代理

func  IsProxy(ctx *gin.Context) bool {
	proxyPath:=[]string{"/k8s/clusters/","/grafana/","/cicd/","/rpc/",
		"/gitlab/", "/nexus/","/harbor/",
	}
	for _,p:=range proxyPath{

		if strings.HasPrefix(ctx.Request.URL.Path,p){
			return true
		}
	}
	return false
}