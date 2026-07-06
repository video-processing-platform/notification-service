package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/alimarzban99/notification-service/internal/bootstrap"
	httpHandler "github.com/alimarzban99/notification-service/internal/interfaces/http"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Logger.Sync()

	http.HandleFunc("/health", httpHandler.HealthCheck)
	http.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.HttpPort),
		Handler: nil,
	}

	errChan := make(chan error, 2)

	// HTTP Server
	go func() {
		app.Logger.Info("HTTP server started",
			zap.Int("port", app.Config.Server.HttpPort),
		)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// gRPC Server
	go func() {
		app.Logger.Info("gRPC server started")

		if err := app.GRPC.Start(); err != nil {
			errChan <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case <-stop:
		app.Logger.Info("shutdown signal received")

	case err := <-errChan:
		app.Logger.Error("server stopped unexpectedly", zap.Error(err))
	}

	// Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.Logger.Error("failed to shutdown http server", zap.Error(err))
	}

	app.GRPC.Stop()

	app.Logger.Info("application stopped")
}
