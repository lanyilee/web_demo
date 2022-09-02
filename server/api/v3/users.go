package v3

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"webase-server/models"
	"webase-server/pkg/errpkg"
	"webase-server/pkg/utils/uuid"
	"webase-server/server/api/base"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	base.BaseV3Controller
}

type createUserParam struct {
	models.LocalUser
	Password string `json:"password"`
}

func (uc *UserController) Create(c *gin.Context) {
	var u createUserParam
	u = createUserParam{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	u.ParentID = uc.GetUserInfo(c).UserID
	u.LocalUser.Password = u.Password

	u.ID = uuid.NewUuidV4String()
	err = uc.APIManager.Store.User.Create(&u.LocalUser)
	if err != nil {
		uc.APIManager.Logger.Info(err)
		uc.SetErrMsg(c, err, 500)
		return
	}

	uc.SetResult(c, nil, 204)
}

func (uc *UserController) Delete(c *gin.Context) {
	userID := c.Param("uid")
	u, isExist, err := uc.APIManager.Store.User.Get(userID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	if u.Username == "admin" {
		uc.SetErrMsg(c, errpkg.APIForbidden)
		return
	}
	if isExist {
		err = uc.APIManager.Store.User.Delete(userID)
		if err != nil {
			uc.Logger.Error("delete user error:", err)
			uc.SetErrMsg(c, err, 500)
			return
		}
	}
	uc.SetResult(c, nil, 204)
}

func (uc *UserController) List(c *gin.Context) {
	_, err := uc.APIManager.Store.User.ListByParentID(uc.GetUserInfo(c).UserID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	var users []models.LocalUser
	uc.SetResult(c, users, 200)
}

func (uc *UserController) Get(c *gin.Context) {
	userID := c.Param("uid")
	u, isExist, err := uc.APIManager.Store.User.Get(userID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	if !isExist {
		uc.SetErrMsg(c, errors.New("not found"), 400)
		return
	}

	uc.SetResult(c, u, 200)
}

func (uc *UserController) Update(c *gin.Context)  {
	userID := c.Param("uid")
	_, isExist, err := uc.APIManager.Store.User.Get(userID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	if !isExist {
		uc.SetErrMsg(c, errors.New("not found"), 400)
		return
	}
	var u createUserParam
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	u.ID = userID
	err = uc.APIManager.Store.User.Update(&u.LocalUser)
	if err != nil {
		uc.APIManager.Logger.Info(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	uc.SetResult(c, nil, 204)
}

func (uc *UserController) Current(c *gin.Context) {
	userID := uc.GetUserInfo(c).UserID
	u, isExist, err := uc.APIManager.Store.User.Get(userID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	if !isExist {
		uc.SetErrMsg(c, errors.New("not found"), 400)
		return
	}
	uc.SetResult(c, u, 200)
}

type putCurrentParam struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (uc *UserController) UpdateCurrent(c *gin.Context) {
	userID := uc.GetUserInfo(c).UserID
	var in putCurrentParam
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		uc.SetErrMsg(c, err, 400)
		return
	}
	u, isExist, err := uc.APIManager.Store.User.Get(userID)
	if err != nil {
		uc.Logger.Error(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	if !isExist {
		uc.SetErrMsg(c, errors.New("not found"), 400)
		return
	}
	if in.NewPassword != "" {
		err = uc.APIManager.Store.User.ChangePassword(u.Username, in.OldPassword, in.NewPassword)
		if err != nil {
			uc.APIManager.Logger.Info(err)
			uc.SetErrMsg(c, err, 500)
			return
		}
		uc.SetResult(c, nil, 204)
		return
	}
	u.Name = in.Name
	u.Email = in.Email
	u.Phone = in.Phone
	err = uc.APIManager.Store.User.Update(u)
	if err != nil {
		uc.APIManager.Logger.Info(err)
		uc.SetErrMsg(c, err, 500)
		return
	}
	uc.SetResult(c, nil, 204)
}
