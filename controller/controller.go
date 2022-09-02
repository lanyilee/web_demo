package controller

import (
	"webase-server/cmd/server/config"
	"webase-server/models"
	"webase-server/pkg/rbac"
)

//Start ...
//func Start(store *models.Store, dialer models.DialerFactory) error {
//	c := newCluster(store, dialer)
//	go func() {
//		for {
//			err := c.watch(context.Background())
//			if err != nil {
//				zap.S().Error(err)
//			}
//		}
//
//	}()
//	return nil
//}

func StartCasbin(store *models.Store, cfg *config.Config) error {
	c := rbac.NewCasbinMysqlStore(store, cfg)
	// c.UseBeegoXormForCasbinClient(cfg)
	c.UseXorm(cfg)
	c.CreateCasbinPolicy(store)
	return nil
}
