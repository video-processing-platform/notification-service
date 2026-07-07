package mail

import (
	"context"
	"fmt"
	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"
	"github.com/alimarzban99/notification-service/internal/infrastructure/metrics"
	"github.com/alimarzban99/notification-service/internal/interfaces/mail"
	"go.uber.org/zap"
	"net"
	"net/smtp"
	"time"
)

type MailHog struct {
	host string
	port int
	from string
	auth smtp.Auth
	log  logger.Logger
}

func NewMailHog(cfg config.MailConfig, log logger.Logger) mail.MailService {
	return &MailHog{
		host: cfg.MailHog.Host,
		port: cfg.MailHog.Port,
		from: cfg.From,
		auth: nil,
		log:  log,
	}
}

func (m *MailHog) Send(
	_ context.Context,
	to string,
	subject string,
	body string,
) error {

	start := time.Now()

	defer func() {
		metrics.EmailSendDuration.Observe(time.Since(start).Seconds())
	}()

	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(
		addr,
		m.auth,
		m.from,
		[]string{to},
		msg,
	)

	if err != nil {
		metrics.EmailsFailedTotal.Inc()

		m.log.Error("mailhog send failed",
			zap.String("to", to),
			zap.String("error", err.Error()),
		)
		return err
	}

	metrics.EmailsSentTotal.Inc()

	m.log.Info("Email sent via mailhog",
		zap.String("to", to),
	)

	return nil
}

func (m *MailHog) Ping(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", m.host, m.port)

	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			m.log.Warn("failed to close smtp connection",
				zap.Error(err),
			)
		}
	}()

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Close(); err != nil {
			m.log.Warn("failed to close smtp client",
				zap.Error(err),
			)
		}
	}()
	return nil
}
