package bootstrap

import (
	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/application/notification"
	grpcserver "github.com/alimarzban99/notification-service/internal/infrastructure/grpc"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"
	"github.com/alimarzban99/notification-service/internal/infrastructure/mail"
	"github.com/alimarzban99/notification-service/internal/infrastructure/metrics"
)

type App struct {
	Config *config.Config
	Logger logger.Logger
	GRPC   *grpcserver.Server
}

func New() (*App, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	metrics.Register()

	mailer := mail.NewMailService(cfg.Mail, log)

	notificationService := notification.NewService(log, mailer)

	grpcServer := grpcserver.New(cfg, log, notificationService)

	app := &App{
		Config: cfg,
		Logger: log,
		GRPC:   grpcServer,
	}

	return app, nil
}
