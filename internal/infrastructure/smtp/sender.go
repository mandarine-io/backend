package smtp

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
	"strings"
)

type Sender interface {
	SendPlainMessage(subject string, content string, from string, to string, attachments ...string) error
	SendPlainMessages(subject string, content string, from string, to []string, attachments ...string) error
	SendHTMLMessage(subject string, content string, from string, to string, attachments ...string) error
	SendHTMLMessages(subject string, content string, from string, to []string, attachments ...string) error
}

type Option func(s *sender) error

func WithLogger(logger zerolog.Logger) Option {
	return func(s *sender) error {
		s.logger = logger
		return nil
	}
}

type sender struct {
	dialer *gomail.Dialer
	logger zerolog.Logger
}

func NewSender(dialer *gomail.Dialer, opts ...Option) (Sender, error) {
	s := &sender{
		dialer: dialer,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return s, nil
}

func (s *sender) SendPlainMessages(
	subject string,
	content string,
	from string,
	to []string,
	attachments ...string,
) error {
	s.logger.Debug().Msgf("sending plain email to %s", strings.Join(to, ","))

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	err := s.dialer.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send plain email: %w", err)
	}

	return err
}

func (s *sender) SendPlainMessage(subject string, content string, from string, to string, attachments ...string) error {
	return s.SendPlainMessages(subject, content, from, []string{to}, attachments...)
}

func (s *sender) SendHTMLMessages(
	subject string,
	content string,
	from string,
	to []string,
	attachments ...string,
) error {
	s.logger.Debug().Msgf("sending HTML email to %s", strings.Join(to, ","))

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	err := s.dialer.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send HTML email: %w", err)
	}

	return err
}

func (s *sender) SendHTMLMessage(subject string, content string, from string, to string, attachments ...string) error {
	return s.SendHTMLMessages(subject, content, from, []string{to}, attachments...)
}
