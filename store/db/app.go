package db

import (
	"webase-server/models"

	"xorm.io/xorm"
)

type app struct {
	db *xorm.Engine
}

//NewApp 应用存储
func NewApp(db *xorm.Engine) models.AppStore {
	return &app{db}
}

func (s *app) Create(h *models.App) error {
	_, err := s.db.Insert(h)
	return err
}

func (s *app) Update(h *models.App) error {
	_, err := s.db.Update(h, &models.App{ID: h.ID})
	return err
}

func (s *app) Delete(id string) error {
	_, err := s.db.Delete(models.App{ID: id})
	return err
}

func (s *app) Find(filter *models.App) ([]models.App, error) {
	var apps []models.App
	err := s.db.Find(&apps, filter)
	if err != nil {
		return apps, err
	}
	return apps, err
}

func (s *app) List() ([]models.App, error) {
	apps := make([]models.App, 0)
	err := s.db.Find(&apps)
	if err != nil {
		return apps, err
	}
	return apps, err
}
