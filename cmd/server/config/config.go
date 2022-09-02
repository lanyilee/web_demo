package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          int    `envconfig:"WEBASE_SERVER_PORT"`
	Host          string `envconfig:"WEBASE_SERVER_HOST"`
	Name          string `envconfig:"WEBASE_SERVER_NAME" default:"webase"`
	EnableHTTPS   bool   `envconfig:"WEBASE_SERVER_ENABLEHTTPS"`
	HTTPSCertFile string `envconfig:"WEBASE_SERVER_HTTPSCERTFILE"`
	HTTPSKeyFile  string `envconfig:"WEBASE_SERVER_HTTPSKEYFILE"`
	Mysql         Mysql
	Debug         bool `envconfig:"WEBASE_SERVER_DEBUG"`
	LogJSON       bool `envconfig:"WEBASE_SERVER_LOG_JSON"`
	Exporter      ExporterOpt
}

type Mysql struct {
	Username string `envconfig:"WEBASE_SERVER_MYSQL_USERNAME"`
	Password string `envconfig:"WEBASE_SERVER_MYSQL_PASSWORD"`
	Address  string `envconfig:"WEBASE_SERVER_MYSQL_ADDRESS"`
	Database string `envconfig:"WEBASE_SERVER_MYSQL_DATABASE"`
}




type ExporterOpt struct {
	ServerEnabled bool   `envconfig:"WEBASE_SERVER_EXPORTER_ENABLED"`
	ServerPort    string `envconfig:"WEBASE_SERVER_EXPORTER_PORT" default:"9121"`
	UIEnabled     bool   `envconfig:"WEBASE_SERVER_UI_EXPORTER_ENABLED"`
	UIPort        string `envconfig:"WEBASE_SERVER_UI_EXPORTER_PORT" default:"9122"`
	// TCPReportEnabled bool `envconfig:"WEBASE_SERVER_TCPREPORT_ENABLED"`
	// TCPReportPort    int  `envconfig:"WEBASE_SERVER_TCPREPORT_PORT" default:"9123"`
}

func Load() (*Config, error) {
	godotenv.Load(".env")
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	//if cfg.Debug {
	//	logger.Setup("debug", cfg.LogJSON)
	//} else {
	//	logger.Setup("info", cfg.LogJSON)
	//}
	return &cfg, err
}
