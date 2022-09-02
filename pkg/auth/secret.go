package auth

import (
	"errors"
	"net/http"
	"strings"
	"webase-server/models"
)

func (c *jwtClient) ParseFromSecret(req *http.Request) (models.UserInfo, error) {
	secret := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	bind, exist, err := c.store.UserBind.GetBindBySecret(secret)
	if err != nil {
		return models.UserInfo{}, err
	}
	if !exist {
		return models.UserInfo{}, errors.New("token not found")
	}
	info := models.UserInfo{
		UserID: bind.UserID,
	}
	return info, nil
}
