package mail

import (
	"context"
	"fmt"
	"github.com/alimarzban99/notification-service/internal/infrastructure/metrics"
	"net"
	"net/smtp"
	"time"

	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"
	"go.uber.org/zap"
)

type MailTrap struct {
	host          string
	port          int
	username      string
	password      string
	from          string
	log           logger.Logger
	timeout       int
	retryAttempts int
	retryDelay    int
}

func NewMailTrap(cfg config.MailConfig, from string, log logger.Logger) *MailTrap {
	return &MailTrap{
		host:          cfg.MailTrap.Host,
		port:          cfg.MailTrap.Port,
		username:      cfg.MailTrap.Username,
		password:      cfg.MailTrap.Password,
		from:          from,
		log:           log,
		timeout:       cfg.Timeout,
		retryAttempts: cfg.RetryAttempts,
		retryDelay:    cfg.RetryDelay,
	}
}

func (m *MailTrap) Send(ctx context.Context, to, subject, body string) error {

	timeout := time.Duration(m.timeout) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

	delay := time.Duration(m.retryDelay) * time.Second
	err := retry(
		ctx,
		m.retryAttempts,
		delay,
		func() error {

			return smtp.SendMail(
				addr,
				auth,
				m.from,
				[]string{to},
				msg,
			)

		},
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

func (m *MailTrap) Ping(ctx context.Context) error {
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
		if err := client.Quit(); err != nil {
			m.log.Warn("failed to quit smtp client",
				zap.Error(err),
			)
		}
	}()

	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return nil
}
