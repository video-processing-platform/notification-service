package service

import (
	"context"
	"fmt"
	"github.com/alimarzban99/notification-service/internal/application/notification"
	domain "github.com/alimarzban99/notification-service/internal/interfaces/mail"

	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"

	"go.uber.org/zap"
)

const (
	ProcessingCompletedSubject = "Processing Completed"
	ProcessingCompletedBody    = "Video %s has been processed."
)

type service struct {
	logger logger.Logger
	mailer domain.MailService
}

func NewService(
	log logger.Logger,
	mailer domain.MailService,
) notification.Service {

	return &service{
		logger: log,
		mailer: mailer,
	}
}

func (s *service) SendEmail(ctx context.Context, input notification.SendEmailInput) error {

	s.logger.Info(
		"sending email",
		zap.String("email", input.Email),
		zap.String("subject", input.Subject),
	)

	if err := s.mailer.Send(
		ctx,
		input.Email,
		input.Subject,
		input.Body,
	); err != nil {

		s.logger.Error(
			"failed to send email",
			zap.String("email", input.Email),
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
		"email sent successfully",
		zap.String("email", input.Email),
	)

	return nil
}

func (s *service) SendProcessingCompleted(ctx context.Context, input notification.ProcessingCompletedInput) error {

	s.logger.Info(
		"sending processing completed email",
		zap.String("email", input.Email),
		zap.String("filename", input.Filename),
	)

	body := fmt.Sprintf(
		ProcessingCompletedBody,
		input.Filename,
	)

	if err := s.mailer.Send(
		ctx,
		input.Email,
		ProcessingCompletedSubject,
		body,
	); err != nil {

		s.logger.Error(
			"failed to send processing completed email",
			zap.String("email", input.Email),
			zap.String("filename", input.Filename),
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
		"processing completed email sent successfully",
		zap.String("email", input.Email),
		zap.String("filename", input.Filename),
	)

	return nil
}
