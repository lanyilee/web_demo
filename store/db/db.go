package db

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"webase-server/cmd/server/config"
	"webase-server/models"
	"webase-server/pkg/logger"
	"xorm.io/core"
	"xorm.io/xorm"
)

func New(cfg *config.Config) (*models.Store, error) {
	//datasource := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ")/?charset=utf8"
	//db, err := sql.Open("mysql", datasource)
	//if err != nil {
	//	return nil, err
	//}
	//defer db.Close()
	//db.Exec("CREATE DATABASE " + cfg.Mysql.Database + " DEFAULT CHARACTER SET utf8 COLLATE utf8_bin")
	//_, err = db.Exec("use " + cfg.Mysql.Database)
	//if err != nil {
	//	return nil, err
	//}
	//if err := mysql.Migrate(db); err != nil {
	//	return nil, err
	//}
	datasource := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ")/" + cfg.Mysql.Database + "?charset=utf8"
	engine, err := xorm.NewEngine("mysql", datasource)
	if err != nil {
		return nil, err
	}
	engine.SetLogger(logger.NewXormZapLogger(zap.L()))
	engine.SetLogLevel(core.LOG_ERR)
	store := &models.Store{
		User:          NewUser(engine),
		App:           NewApp(engine),
		UserBind:      NewUser(engine),
	}
	return store, nil
}
