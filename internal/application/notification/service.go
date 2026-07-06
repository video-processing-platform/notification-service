package notification

import "context"

type Service interface {
	SendEmail(
		ctx context.Context,
		input SendEmailInput,
	) error

	SendProcessingCompleted(
		ctx context.Context,
		input ProcessingCompletedInput,
	) error
}
