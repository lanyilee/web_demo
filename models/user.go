package models

import "time"

type LocalUser struct {
	ID       string    `xorm:"id" json:"id"`
	ParentID string    `xorm:"parent_id" json:"parent_id"`
	Username string    `xorm:"username" json:"username"`
	Password string    `xorm:"password" json:"-"`
	Name     string    `xorm:"name" json:"name"`
	Role     string    `xorm:"role" json:"role"`
	Level    int       `xorm:"level" json:"level"`
	Phone    string    `xorm:"phone" json:"phone"`
	Email    string    `xorm:"email" json:"email"`
	Created  time.Time `xorm:"created" json:"created"`
	Updated  time.Time `xorm:"updated" json:"updated"`
}

func (u *LocalUser) TableName() string {
	return "local_user"
}

type LocalUserBind struct {
	ID         string `xorm:"id" json:"bind_id"`
	UserID     string `xorm:"user_id" json:"user_id"`
	Provider   string `xorm:"provider" json:"provider"`
	ProviderID string `xorm:"provider_id" json:"provider_id"`
}

type LocalUserBindExtends struct {
	LocalUserBind `xorm:"extends"`
	LocalUser     `xorm:"extends"`
}

type UserStore interface {
	Create(u *LocalUser) error
	Update(u *LocalUser) error
	Delete(id string) error
	List() ([]LocalUser, error)
	ListByParentID(pid string) ([]LocalUser, error)
	Get(id string) (*LocalUser, bool, error)
	GetUserByUsername(name string) (*LocalUser, bool, error)
	Login(username, password string) (*LocalUser, error)
	ChangePassword(username, oldPassword, newPassword string) error
}

type UserBindStore interface {
	Bind(bind *LocalUserBind) error
	DeleteBind(userID, provider string) error
	FindBindsByUser(id string) ([]LocalUserBind, error)
	GetBindBySecret(secret string) (*LocalUserBind, bool, error)
}
