package grpcserver

import (
	"context"
	"github.com/alimarzban99/notification-service/internal/application/notification"

	notificationpb "github.com/alimarzban99/notification-service/proto/notification"
)

type Handler struct {
	notificationpb.UnimplementedNotificationServiceServer

	service notification.Service
}

func NewHandler(
	service notification.Service,
) *Handler {

	return &Handler{
		service: service,
	}
}

func (h *Handler) SendEmail(ctx context.Context, req *notificationpb.SendEmailRequest) (*notificationpb.SendEmailResponse, error) {

	err := h.service.SendEmail(
		ctx,
		notification.SendEmailInput{
			Email:   req.Email,
			Subject: req.Subject,
			Body:    req.Body,
		},
	)

	if err != nil {
		return nil, err
	}

	return &notificationpb.SendEmailResponse{
		Success: true,
		Message: "Email queued",
	}, nil
}

func (h *Handler) SendProcessingCompleted(ctx context.Context, req *notificationpb.ProcessingCompletedRequest) (*notificationpb.ProcessingCompletedResponse, error) {

	err := h.service.SendProcessingCompleted(
		ctx,
		notification.ProcessingCompletedInput{
			VideoID:  req.VideoId,
			Email:    req.Email,
			Filename: req.Filename,
			Status:   req.Status,
		},
	)

	if err != nil {
		return nil, err
	}

	return &notificationpb.ProcessingCompletedResponse{
		Success: true,
		Message: "Notification queued",
	}, nil

}
