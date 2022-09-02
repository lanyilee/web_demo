package rbac

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"webase-server/cmd/server/config"
	"webase-server/models"

	// beegoormadapter "github.com/casbin/beego-orm-adapter/v3"

	casbin "github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
)

type casbinMysql struct {
	store        *models.Store
	cfg          *config.Config
	casbinPolicy *CasbinPolicy
}

func NewCasbinMysqlStore(store *models.Store, cfg *config.Config) *casbinMysql {
	casbinPolicy := &CasbinPolicy{store}
	return &casbinMysql{store, cfg, casbinPolicy}
}

func finalizer(db *sql.DB) {
	err := db.Close()
	if err != nil {
		zap.L().Error(err.Error())
	}
}

var Enf *casbin.SyncedEnforcer

// func (s *casbinMysql) UseMysql(cfg *config.Config) (*casbin.Enforcer, error) {
// 	var err error
// 	if !(Enf == nil) {
// 		zap.S().Info("I am already init")
// 		return Enf, nil
// 	}
// 	zap.S().Info("lazy init")
// 	// connect to the database first.
// 	datasource := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ")/" + cfg.Mysql.Database + "?charset=utf8"
// 	// db, err := sql.Open("mysql", "root:root123@tcp(188.8.5.1:33060)/webase")
// 	db, err := sql.Open("mysql", datasource)
// 	if err != nil {
// 		log.Println("open mysql error is: ", err)
// 		return nil, err
// 	}
// 	if err = db.Ping(); err != nil {
// 		log.Println("ping mysql error is: ", err)
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(20)
// 	db.SetMaxIdleConns(10)
// 	db.SetConnMaxLifetime(time.Minute * 10)

// 	// need to control by user, not the package
// 	runtime.SetFinalizer(db, finalizer)

// 	// Initialize an adapter and use it in a Casbin enforcer:
// 	// The adapter will use the Sqlite3 table name "casbin_rule_test",
// 	// the default table name is "casbin_rule".
// 	// If it doesn't exist, the adapter will create it automatically.
// 	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
// 	if err != nil {
// 		log.Println("sql adapter error is: ", err)
// 		return nil, err
// 	}
// 	casbin_model := cfg.Casbin.Model
// 	Enf, err = casbin.NewEnforcer(casbin_model, a)
// 	if err != nil {
// 		log.Println("do not meet the definition, error is: ", err)
// 		return nil, err
// 	}

// 	// Load the policy from DB.
// 	if err = Enf.LoadPolicy(); err != nil {
// 		log.Println("LoadPolicy failed, err: ", err)
// 		return nil, err
// 	}

// 	// Save the policy back to DB.
// 	// if err = Enf.SavePolicy(); err != nil {
// 	// 	log.Println("SavePolicy failed, err: ", err)
// 	// }
// 	return Enf, nil
// }

// func (s *casbinMysql) UseBeegoXormForCasbinClient(cfg *config.Config) (*casbin.Enforcer, error) {
// 	var err error
// 	if !(Enf == nil) {
// 		zap.S().Info("I am already init")
// 		return Enf, nil
// 	}
// 	zap.S().Info("lazy init")
// 	// datasource := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ")/casbin?charset=utf8"
// 	// a, err := beegoormadapter.NewAdapter("default", "mysql", "root:root123@tcp(188.8.5.1:33060)/webase")
// 	// a, err := beegoormadapter.NewAdapter("default", "mysql", datasource)
// 	// if err != nil {
// 	// 	zap.S().Error(err)
// 	// }
// 	// file_map := GetCasbinFile()
// 	// Enf, err = casbin.NewEnforcer(file_map["model"], a)
// 	casbin_model := cfg.Casbin.Model
// 	casbin_policy := cfg.Casbin.Policy
// 	Enf, err = casbin.NewEnforcer(casbin_model, casbin_policy)
// 	if err != nil {
// 		zap.S().Error(err)
// 	}
// 	Enf.EnableAutoSave(true)
// 	return Enf, nil
// }

func (s *casbinMysql) UseXorm(cfg *config.Config) (*casbin.SyncedEnforcer, error) {
	var err error
	if !(Enf == nil) {
		// fmt.Println("I am already init")
		return Enf, nil
	}
	fmt.Println("lazy init")

	datasource := cfg.Mysql.Username + ":" + cfg.Mysql.Password + "@tcp(" + cfg.Mysql.Address + ")/" + cfg.Mysql.Database + "?charset=utf8"
	zap.S().Info("datasource", datasource)
	//a, _ := xormadapter.NewAdapter("mysql", datasource, true)
	//if err != nil {
	//	zap.S().Error("errors")
	//	return nil, err
	//}
	//casbin_model := cfg.Casbin.Model
	//zap.S().Info("casbin_model", casbin_model)
	//Enf, _ = casbin.NewSyncedEnforcer(casbin_model, a)
	//Enf.EnableAutoSave(true)

	if err = Enf.LoadPolicy(); err != nil {
		zap.S().Info("LoadPolicy failed, err: ", err)
		return nil, err
	}
	return Enf, nil
}

func (s *casbinMysql) CreateCasbinPolicy(store *models.Store) error {
	e, err := s.UseXorm(s.cfg)
	if err != nil {
		zap.S().Error("casbin enforce errors:", err)
		return err
	}
	e.AddPolicy("role:admin", "*", "*")
	e.AddGroupingPolicy("admin", "role:admin")
	// e.SavePolicy()
	s.casbinPolicy.ClusterPolicy(store)
	s.casbinPolicy.NamespacePolicy(store)
	s.casbinPolicy.ProjectPolicy(store)
	s.casbinPolicy.BindingPolicy(store)
	return nil
}
