package mail

import (
	"context"
	"fmt"
	"github.com/alimarzban99/notification-service/internal/infrastructure/metrics"
	"net/smtp"
	"time"

	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"
	"go.uber.org/zap"
)

type MailTrap struct {
	host     string
	port     int
	username string
	password string
	from     string
	log      logger.Logger
}

func NewMailTrap(cfg config.MailTrapConfig, from string, log logger.Logger) *MailTrap {
	return &MailTrap{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     from,
		log:      log,
	}
}

func (m *MailTrap) Send(ctx context.Context, to, subject, body string) error {

	start := time.Now()
	defer func() {
		metrics.EmailSendDuration.Observe(time.Since(start).Seconds())
	}()

	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	auth := smtp.PlainAuth(
		"",
		m.username,
		m.password,
		m.host,
	)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(
		addr,
		auth,
		m.from,
		[]string{to},
		msg,
	)

	if err != nil {
		metrics.EmailsFailedTotal.Inc()

		m.log.Error("Mailtrap send failed",
			zap.String("error", err.Error()),
		)
		return err
	}

	metrics.EmailsSentTotal.Inc()
	m.log.Info("Mail sent via Mailtrap",
		zap.String("to", to),
	)

	return nil
}
