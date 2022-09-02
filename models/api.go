package models

import (
	"crypto/tls"
	"net/http"
	"time"
	"webase-server/cmd/server/config"
)

type APIManager struct {
	Port          int
	Logger        Logger
	Config        *config.Config
	Store         *Store
	DockerCompose AppService
	Auth          AuthInfterface
	//OIDC          OIDCInterface
}

var DefaultHTTPClient = &http.Client{
	Timeout: 30 * time.Minute,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type User struct {
	UserRole    string `json:"userRole"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	DroneToken  string `json:"drone_token"`
	ID          string `json:"id"`
}


type CreateTypes struct {
	Namespace string `json:"namespace"`
}
type Pagination struct {
	Limit int `json:"type"`
	Total int `json:"total"`
}


