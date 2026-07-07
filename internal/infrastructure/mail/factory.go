package mail

import (
	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"
	"github.com/alimarzban99/notification-service/internal/interfaces/mail"
)

func NewMailService(cfg config.MailConfig, log logger.Logger) mail.MailService {

	switch cfg.Provider {
	case "mailtrap":
		return NewMailTrap(cfg, cfg.From, log)
	case "mailhog":
		fallthrough
	default:
		return NewMailHog(cfg, log)
	}
}
