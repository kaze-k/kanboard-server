package router

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"server/api/admin"
	"server/api/kanboard"
	"server/internal/constant"
	"server/internal/global"
	"server/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := getRouter()
	initRoutes(router)
	initServer(
		constant.ServerConfig.Host,
		constant.ServerConfig.Port,
		constant.ServerConfig.Name,
		router,
	)
}

func initRoutes(router *gin.Engine) {
	router.Static(constant.FileConfig.Static, constant.FileConfig.Path)
	api := router.Group("/api")
	{
		admin.InitApi(api)
		kanboard.InitApi(api)
	}
}

func getRouter() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.Cors(),
	)

	return router
}

func initServer(host string, port string, name string, handler http.Handler) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		global.Logger.Infof("%s server start at %s", name, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error(err)
			panic(err)
		}
	}()

	<-ctx.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), constant.ServerConfig.Timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Error(err)
		panic(err)
	}
}
