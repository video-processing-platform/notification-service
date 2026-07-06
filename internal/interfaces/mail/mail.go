package mail

import (
	"context"
)

type MailService interface {
	Send(
		ctx context.Context,
		to string,
		subject string,
		body string,
	) error
}
