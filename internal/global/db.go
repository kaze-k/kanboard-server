package global

import (
	"io"
	"log"
	"os"

	"server/internal/constant"
	"server/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDB() (*gorm.DB, error) {
	logMode := logger.Error
	colorful := false
	if constant.EnvConfig.Mode == "debug" {
		logMode = logger.Info
		colorful = true
	}

	db, err := gorm.Open(mysql.Open(constant.DBConfig.Dsn), &gorm.Config{
		Logger: initDBLogger(logMode, colorful),
	})
	if err != nil {
		return nil, err
	}

	autoMigrate(db)

	connPoolErr := initDBConnPool(db)
	if connPoolErr != nil {
		return nil, connPoolErr
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Project{},
		&models.ProjectMember{},
		&models.TaskAssignee{},
		&models.Resource{},
	)
	if err != nil {
		Logger.Error(err)
		panic(err)
	}
}

func initDBLogger(level logger.LogLevel, colorful bool) logger.Interface {
	return logger.New(
		log.New(io.MultiWriter(GetWriter(), os.Stdout), "\n", log.LstdFlags),
		logger.Config{
			LogLevel: level,
			Colorful: colorful,
		},
	)
}

func initDBConnPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(constant.DBConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(constant.DBConfig.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(constant.DBConfig.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(constant.DBConfig.ConnMaxLifetime)

	return nil
}
