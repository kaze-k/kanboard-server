package config

import (
	"io"
	"os"
	"server/internal/constant"
	"server/internal/global"

	"github.com/gin-gonic/gin"
)

func initGinConfig() {
	if constant.EnvConfig.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if constant.EnvConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	gin.DefaultErrorWriter = io.MultiWriter(global.GetWriter(), os.Stdout)
}
