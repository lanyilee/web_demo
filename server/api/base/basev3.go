package base

import (
	"context"
	"fmt"
	"webase-server/models"
	"webase-server/pkg/errpkg"

	"github.com/gin-gonic/gin"
)

type BaseV3Controller struct {
	Logger     models.Logger
	APIManager *models.APIManager
}

func (b *BaseV3Controller) GetUserContext(c *gin.Context) context.Context {
	info := c.Param("user_info")
	return context.WithValue(context.Background(), "user_info", info)
}

func (b *BaseV3Controller) GetUserInfo(c *gin.Context) models.UserInfo {
	userinfo, ok := c.Get("user_info")
	if ok == false {
		fmt.Println("解析user_info出错")
	} else {
		fmt.Println(userinfo)
	}
	return userinfo.(models.UserInfo)
	//return c.Ctx.Input.GetData("user_info").(models.UserInfo)
}

func (c *BaseV3Controller) IsAdmin(cg *gin.Context) bool {
	//info := cg.Value("user_info").(models.UserInfo)
	//return info.Role == models.AdminRole
	return true
}

func (c *BaseV3Controller) CheckPermission(cg *gin.Context, roles ...string) bool {
	info := cg.Value("user_info").(models.UserInfo)
	for _, role := range roles {
		if info.Role == role {
			return true
		}
	}
	return false
}

// SetErrMsg return err msg to http
func (b *BaseV3Controller) SetErrMsg(c *gin.Context, err error, codes ...int) {
	statusCode := 500
	if len(codes) > 0 {
		statusCode = codes[0]
	}
	errCode := statusCode
	if len(codes) > 1 {
		errCode = codes[1]
	}
	if msg, ok := err.(errpkg.ErrorMsg); ok {
		b.SetResult(c, msg, statusCode)
		return
	}
	b.SetResult(c, errpkg.NewAPIError(errCode, err), statusCode)
}

//SetResult return json to http
func (b *BaseV3Controller) SetResult(c *gin.Context, result interface{}, code int, key ...string) {
	//c.Status(code)
	if result == nil && (len(key) == 0 || key[0] == "") {
		return
	}
	// if len(key) == 0 || key[0] == "" {
	// 	//c.Data["json"] = result
	// } else {
	// 	c.Data["json"] = map[string]interface{}{key[0]: result}
	// }
	c.JSON(code, result)
	//c.ServeJSON()
}
