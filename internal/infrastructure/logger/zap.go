package logger

import (
	"strings"

	"github.com/alimarzban99/notification-service/config"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New(cfg config.LoggerConfig) (*ZapLogger, error) {

	var zapConfig zap.Config

	if strings.ToLower(cfg.Level) == "debug" {

		zapConfig = zap.NewDevelopmentConfig()

	} else {

		zapConfig = zap.NewProductionConfig()

	}

	logger, err := zapConfig.Build()

	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger,
	}, nil
}

func (z *ZapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}
