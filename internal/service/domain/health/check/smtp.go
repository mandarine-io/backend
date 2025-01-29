package check

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
)

type SMTPCheck struct {
	dialer *gomail.Dialer
	logger zerolog.Logger
}

type SMTPOption func(c *SMTPCheck) error

func WithSMTPLogger(logger zerolog.Logger) SMTPOption {
	return func(c *SMTPCheck) error {
		c.logger = logger
		return nil
	}
}

func NewSMTPCheck(dialer *gomail.Dialer, opts ...SMTPOption) (*SMTPCheck, error) {
	check := &SMTPCheck{
		dialer: dialer,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(check); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return check, nil
}

func (c *SMTPCheck) Pass() bool {
	c.logger.Debug().Msgf("check connection to smtp server")

	closer, err := c.dialer.Dial()
	if err != nil {
		c.logger.Error().Stack().Err(err).Msg("failed to connect to smtp server")
		return false
	}

	if err := closer.Close(); err != nil {
		c.logger.Error().Stack().Err(err).Msg("failed to close connection to smtp server")
		return false
	}

	return true
}

func (c *SMTPCheck) Name() string {
	return "smtp"
}
