package db

import (
	"errors"
	"webase-server/models"

	"webase-server/pkg/utils/crypto"

	"xorm.io/xorm"
)

type userStore struct {
	db *xorm.Engine
}

func NewUser(db *xorm.Engine) *userStore {
	return &userStore{db}
}

func (s *userStore) Create(u *models.LocalUser) error {
	u.Password = crypto.MD5(u.Password)
	_, isExist, err := s.GetUserByUsername(u.Username)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("用户已存在")
	}
	_, err = s.db.Insert(u)
	return err
}

func (s *userStore) Update(u *models.LocalUser) error {
	if u.Password != "" {
		u.Password = crypto.MD5(u.Password)
	}
	_, err := s.db.Where("id=?", u.ID).Update(u)
	return err
}

func (s *userStore) Delete(id string) error {
	var u models.LocalUser
	_, err := s.db.Where("id=?", id).Delete(&u)
	return err
}

func (s *userStore) List() ([]models.LocalUser, error) {
	var us []models.LocalUser
	err := s.db.Find(&us)
	return us, err
}

func (s *userStore) ListByParentID(pid string) ([]models.LocalUser, error) {
	var us []models.LocalUser
	err := s.db.Where("parent_id=?", pid).Find(&us)
	return us, err
}

func (s *userStore) Get(id string) (*models.LocalUser, bool, error) {
	var u models.LocalUser
	isExist, err := s.db.Where("id=?", id).Get(&u)
	return &u, isExist, err
}

func (s *userStore) GetUserByUsername(username string) (*models.LocalUser, bool, error) {
	var u models.LocalUser
	isExist, err := s.db.Where("username=?", username).Get(&u)
	return &u, isExist, err
}

func (s *userStore) Login(username, password string) (*models.LocalUser, error) {
	var u models.LocalUser
	isExist, err := s.db.Where("username=? AND password=?", username, crypto.MD5(password)).Get(&u)
	if !isExist {
		return nil, errors.New("用户不存在或密码不正确")
	}
	return &u, err
}

func (s *userStore) ChangePassword(username, oldPassword, newPassword string) error {
	var u models.LocalUser
	isExist, err := s.db.Where("username=? AND password=?", username, crypto.MD5(u.Password)).Get(&u)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("密码不正确")
	}
	return s.Update(&models.LocalUser{ID: u.ID, Password: newPassword})
}

func (s *userStore) Bind(bind *models.LocalUserBind) error {
	var oldBind models.LocalUserBind
	isExist, err := s.db.Where("user_id=? and provider=?", bind.UserID, bind.Provider).Get(&oldBind)
	if err != nil {
		return err
	}
	if isExist {
		bind.ID = oldBind.ID
		_, err := s.db.Where("user_id=? and provider=?", bind.UserID, bind.Provider).Update(bind)
		return err
	}
	_, err = s.db.Insert(bind)
	return err
}

func (s *userStore) DeleteBind(userID, provider string) error {
	var bind models.LocalUserBind
	_, err := s.db.Where("user_id=? and provider=?", userID, provider).Delete(&bind)
	return err
}

func (s *userStore) FindBindsByUser(id string) ([]models.LocalUserBind, error) {
	var bind []models.LocalUserBind
	err := s.db.Where("user_id=?", id).Find(&bind)
	if err != nil {
		return nil, err
	}
	return bind, nil
}

func (s *userStore) GetBindBySecret(secret string) (*models.LocalUserBind, bool, error) {
	var bind models.LocalUserBind
	isExist, err := s.db.Where("provider_id=?", secret).Get(&bind)
	if err != nil {
		return nil, false, err
	}
	return &bind, isExist, nil
}
