package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger *zap.SugaredLogger

	DB *gorm.DB

	Redis *RedisClient
)

func InitGlobal() {
	initLogger()
	initDB()
	initRedis()
}

func initLogger() {
	Logger = getLogger()
}

func initDB() {
	db, err := getDB()
	if err != nil {
		Logger.Error(err)
		panic(err)
	}
	DB = db
	Logger.Info("db connected")
}

func initRedis() {
	redis, err := getRedis()
	if err != nil {
		Logger.Error(err)
		panic(err)
	}
	Redis = redis
	Logger.Info("redis connected")
}
