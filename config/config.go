package config

import "time"

var (
	App      *appConfig
	Database *databaseConfig
	Redis    *redisConfig
)

type appConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Log  struct {
		FilePath         string `mapstructure:"path"`
		FileMaxSize      int    `mapstructure:"max_size"`
		BackUpFileMaxAge int    `mapstructure:"back_up_max_age"`
	}
	Pagination struct {
		DefaultSize int `mapstructure:"default_size"`
		MaxSize     int `mapstructure:"max_size"`
	}
}
type databaseConfig struct {
	Master DbConnectOption `mapstructure:"master"`
	Slave  DbConnectOption `mapstructure:"slave"`
}
type DbConnectOption struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen"`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}

type redisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	DB       int    `mapstructure:"db"`
}
