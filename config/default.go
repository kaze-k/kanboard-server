package config

import (
	"time"

	"github.com/spf13/viper"
)

func setEnvDefaultConfig() {
	viper.SetDefault("env.mode", "release")
}

func setServerDefaultConfig() {
	viper.SetDefault("server.name", "Kanboard")
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "9999")
	viper.SetDefault("server.timeout", 5)
}

func setLogDefaultConfig() {
	viper.SetDefault("log.path", "./logs/")
	viper.SetDefault("log.maxSize", 500)
	viper.SetDefault("log.maxBackups", 10)
	viper.SetDefault("log.maxAge", 30)
	viper.SetDefault("log.compress", true)
}

func setDBDefaultConfig() {
	viper.SetDefault("db.maxIdleConns", 10)
	viper.SetDefault("db.maxOpenConns", 10)
	viper.SetDefault("db.connMaxIdleTime", time.Hour)
	viper.SetDefault("db.connMaxLifetime", time.Hour)
}

func setRedisDefaultConfig() {
	viper.SetDefault("redis.addr", "127.0.0.1:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
}

func setJWTDefaultConfig() {
	viper.SetDefault("jwt.kanboardTokenExpiration", 30*24*time.Hour)
	viper.SetDefault("jwt.adminTokenExpiration", 4*time.Hour)
}

func setFileDefaultConfig() {
	viper.SetDefault("file.path", "./files/")
	viper.SetDefault("file.static", "resources")
}
