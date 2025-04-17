package cmd

import (
	"server/config"
	"server/internal/global"
	"server/internal/router"
)

func init() {
	config.InitConfig()
	global.InitGlobal()
}

func Start() {
	router.Run()
}
