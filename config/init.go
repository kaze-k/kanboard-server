package config

import (
	"fmt"
	"time"

	"server/internal/constant"
	"server/internal/types"

	"github.com/spf13/viper"
)

func InitConfig() {
	loadConfig()
	initEnvConfig()
	initServerConfig()
	initLogConfig()
	initDBConfig()
	initRedisConfig()
	initJWTConfig()
	initFileConfig()
	initGinConfig()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%s\n\n", err)
	} else {
		fmt.Printf("config file loaded: %s\n\n", viper.ConfigFileUsed())
	}
}

func initEnvConfig() {
	setEnvDefaultConfig()
	constant.EnvConfig = &types.Env{
		Mode: viper.GetString("env.mode"),
	}
}

func initServerConfig() {
	setServerDefaultConfig()
	constant.ServerConfig = &types.Server{
		Name:    viper.GetString("server.name"),
		Host:    viper.GetString("server.host"),
		Port:    viper.GetString("server.port"),
		Timeout: viper.GetDuration("server.timeout") * time.Second,
	}
}

func initLogConfig() {
	setLogDefaultConfig()
	constant.LogConfig = &types.Log{
		Path:       viper.GetString("log.path"),
		MaxSize:    viper.GetInt("log.maxSize"),
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     viper.GetInt("log.maxAge"),
		Compress:   viper.GetBool("log.compress"),
	}
}

func initDBConfig() {
	setDBDefaultConfig()
	constant.DBConfig = &types.DB{
		Dsn:             viper.GetString("db.dsn"),
		MaxIdleConns:    viper.GetInt("db.maxIdleConns"),
		MaxOpenConns:    viper.GetInt("db.maxOpenConns"),
		ConnMaxLifetime: viper.GetDuration("db.connMaxLifetime") * time.Hour,
	}
}

func initRedisConfig() {
	setRedisDefaultConfig()
	constant.RedisConfig = &types.Redis{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
}

func initJWTConfig() {
	setJWTDefaultConfig()
	constant.JWTConfig = &types.JWT{
		Secret:                  viper.GetString("jwt.secret"),
		KanboardTokenExpiration: viper.GetDuration("jwt.kanboardTokenExpiration") * 24 * time.Hour,
		AdminTokenExpiration:    viper.GetDuration("jwt.adminTokenExpiration") * time.Hour,
	}
}

func initFileConfig() {
	setFileDefaultConfig()
	constant.FileConfig = &types.File{
		Path:   viper.GetString("file.path"),
		Static: viper.GetString("file.static"),
	}
}
