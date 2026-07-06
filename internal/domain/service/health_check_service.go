package service

import (
	"encoding/json"
	"github.com/alimarzban99/notification-service/config"
	grpcserver "github.com/alimarzban99/notification-service/internal/infrastructure/grpc"
	"github.com/alimarzban99/notification-service/internal/interfaces/mail"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status       string            `json:"status"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Timestamp    time.Time         `json:"timestamp"`
	Uptime       string            `json:"uptime"`
	Dependencies map[string]string `json:"dependencies"`
}

var startedAt = time.Now()

func HealthCheck(
	grpcServer *grpcserver.Server,
	smtpClient mail.MailService,
	config *config.Config,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		status := http.StatusOK
		overallStatus := "UP"

		dependencies := map[string]string{
			"smtp":   "UP",
			"grpc":   "UP",
			"config": "UP",
		}

		// gRPC
		if grpcServer == nil {
			dependencies["grpc"] = "DOWN"
			overallStatus = "DOWN"
			status = http.StatusServiceUnavailable
		}

		// SMTP
		if err := smtpClient.Ping(r.Context()); err != nil {
			dependencies["smtp"] = "DOWN"
			overallStatus = "DOWN"
			status = http.StatusServiceUnavailable
		}

		// Config
		if config == nil {
			dependencies["config"] = "DOWN"
			overallStatus = "DOWN"
			status = http.StatusServiceUnavailable
		}

		resp := HealthResponse{
			Status:       overallStatus,
			Service:      "notification-service",
			Version:      "1.0.0",
			Timestamp:    time.Now().UTC(),
			Uptime:       time.Since(startedAt).Round(time.Second).String(),
			Dependencies: dependencies,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
