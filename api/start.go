package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sys/unix"

	"github.com/buonotti/apisense/api/controllers"
	"github.com/buonotti/apisense/api/middleware"
	"github.com/buonotti/apisense/docs"
	"github.com/buonotti/apisense/errors"
	"github.com/buonotti/apisense/log"
)

func Start() error {
	// TODO config
	docs.SwaggerInfo.BasePath = "/api"
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.CORS())
	router.Use(log.GinLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.Limiter())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	api.GET("/health", controllers.GetHealth)
	api.GET("/reports", controllers.AllReports)
	api.GET("/reports/:id", controllers.Report)
	api.GET("/ws", controllers.Ws)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, unix.SIGINT, unix.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.ApiLogger.Error(err.Error())
		}
	}()

	log.ApiLogger.Info("api service started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.ApiLogger.Info("stopping api service")

	if err := srv.Shutdown(ctx); err != nil {
		err = errors.CannotStopApiServiceError.Wrap(err, "cannot stop api service")
	}

	return nil
}
