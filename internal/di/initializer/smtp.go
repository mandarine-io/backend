package initializer

import (
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/smtp"
)

func SMTP(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup SMTP")

		var err error

		smtpConfig := toSMTPConfig(c.Config.SMTP)
		c.Infrastructure.SMTPDialer, err = smtp.NewDialer(smtpConfig)
		if err != nil {
			return err
		}

		c.Infrastructure.SMTPSender, err = smtp.NewSender(
			c.Infrastructure.SMTPDialer,
			smtp.WithLogger(c.Logger.With().Str("component", "smtp-sender").Logger()),
		)

		return err
	}
}

func toSMTPConfig(cfg config.SMTPConfig) smtp.Config {
	return smtp.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		SSL:      cfg.SSL,
	}
}
