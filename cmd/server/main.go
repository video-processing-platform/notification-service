package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/alimarzban99/notification-service/internal/bootstrap"
	"github.com/alimarzban99/notification-service/internal/domain/service"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = app.Logger.Sync()
	}()

	//------------------------------------
	// Gin Router
	//------------------------------------

	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/health", gin.WrapF(
		service.HealthCheck(app.GRPC, app.Mailer, app.Config),
	))

	router.GET("/metrics", gin.WrapH(
		promhttp.Handler(),
	))

	//------------------------------------
	// HTTP Server
	//------------------------------------

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.HTTPPort),
		Handler: router,
	}

	errChan := make(chan error, 2)

	//------------------------------------
	// Start HTTP
	//------------------------------------

	go func() {
		app.Logger.Info(
			"HTTP server started",
			zap.Int("port", app.Config.Server.HTTPPort),
		)

		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			errChan <- err
		}
	}()

	//------------------------------------
	// Start gRPC
	//------------------------------------

	go func() {
		app.Logger.Info("gRPC server started")

		if err := app.GRPC.Start(); err != nil {
			errChan <- err
		}
	}()

	//------------------------------------
	// Wait Signal
	//------------------------------------

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {

	case <-stop:
		app.Logger.Info("shutdown signal received")

	case err := <-errChan:
		app.Logger.Error(
			"server stopped unexpectedly",
			zap.Error(err),
		)
	}

	//------------------------------------
	// Graceful Shutdown
	//------------------------------------

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	app.Logger.Info("stopping HTTP server...")

	if err := httpServer.Shutdown(ctx); err != nil {
		app.Logger.Error(
			"failed to shutdown http server",
			zap.Error(err),
		)
	}

	app.Logger.Info("stopping gRPC server...")

	app.GRPC.Stop()

	app.Logger.Info("application stopped")
}
