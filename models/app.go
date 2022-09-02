package models

import (
	"time"
)

//App 应用
type App struct {
	ID        string                 `json:"app_id" xorm:"id"`
	UserID    string                 `json:"app_user_id" xorm:"user_id"`
	Name      string                 `json:"app_name" xorm:"name"`
	AliasName string                 `json:"app_alias_name" xorm:"alias_name"`
	Kind      string                 `json:"app_kind" xorm:"kind"`
	Content   map[string]interface{} `json:"app_content" xorm:"content"`
	Created   time.Time              `json:"app_created" xorm:"created" `
	Updated   time.Time              `json:"app_updated" xorm:"updated"`
}

//TableName ...
func (app *App) TableName() string {
	return "app"
}


const (
	//DockerComposeAppKind docker compose app
	DockerComposeAppKind = "docker-compose"
)

//AppContent 应用内容
type AppContent interface {
	Hosts() []string
}

//AppStore 应用数据保存
type AppStore interface {
	Create(app *App) error
	Update(app *App) error
	Delete(id string) error
	Find(filter *App) ([]App, error)
	List() ([]App, error)
}

//AppActionOpt 应用操作
type AppActionOpt struct {
	Action   string   `json:"action"`
	Services []string `json:"services"`
	Force    bool     `json:"force"`
}

//AppLogOpt 应用日志
type AppLogOpt struct {
	Services []string `json:"services"`
	Tails    int      `json:"tail"`
}

//AppService 应用部署服务
type AppService interface {
	Create(app *App) error
	Action(id string, opt AppActionOpt) error
	List(user string) ([]App, error)
	Upgrade(app *App) error
	Delete(id string) error
	Get(id string) (*App, error)
	Status(id string) (interface{}, error)
	Logs(id string, opt AppLogOpt) (string, error)
}
