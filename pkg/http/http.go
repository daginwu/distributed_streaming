package http

import (
	"go.uber.org/fx"

	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Modual = fx.Options(

	fx.Provide(
		NewHTTPServer,
	),
	fx.Invoke(
		InitHTTPServer,
	),
)

func NewHTTPServer(lc fx.Lifecycle) *gin.Engine {

	log.Println("[Distributed_Streaming] HTTP server start")

	port := viper.GetString("http.port")
	router := gin.Default()
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return router
}

func InitHTTPServer(router *gin.Engine) error {
	return nil
}
