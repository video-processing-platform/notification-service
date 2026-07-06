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

func (s *MailHog) Send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {

	start := time.Now()

	defer func() {
		metrics.EmailSendDuration.Observe(time.Since(start).Seconds())
	}()

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(
		addr,
		s.auth,
		s.from,
		[]string{to},
		msg,
	)

	if err != nil {
		metrics.EmailsFailedTotal.Inc()

		s.log.Error("mailhog send failed",
			zap.String("to", to),
			zap.String("error", err.Error()),
		)
		return err
	}

	metrics.EmailsSentTotal.Inc()

	s.log.Info("Email sent via mailhog",
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
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}
	defer client.Close()

	return nil
}
