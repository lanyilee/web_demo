package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
	"webase-server/cmd/server/config"
	"webase-server/models"
	"webase-server/pkg/auth"
	"webase-server/pkg/logger"
	"webase-server/pkg/version"
	"webase-server/server"
	"webase-server/store/db"
)

func main() {
	fmt.Println(version.Get())
	cfg, err := config.Load()
	if err != nil {
		zap.S().Error(err)
		os.Exit(-1)
	}
	if cfg.Debug {
		logger.Setup("debug", cfg.LogJSON)
	} else {
		logger.Setup("info", cfg.LogJSON)
	}
	defer func() {
		zap.L().Sync()
		zap.S().Sync()
	}()
	store, err := db.New(cfg)
	if err != nil {
		zap.S().Error(err)
		os.Exit(-1)
	}
	//
	authJwt := auth.New("", 8*time.Hour, store)

	apimgr := &models.APIManager{
		Port:          cfg.Port,
		Logger:        zap.S(), //使用zap
		Config:        cfg,
		//Store:         store,
		Auth:          authJwt,
	}

	zap.S().Info("start")

	server.StartGin(apimgr)
}
