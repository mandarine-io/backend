package smtp

import (
	"crypto/tls"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"strings"
)

type Sender interface {
	HealthCheck() bool
	SendPlainMessage(subject string, content string, to string, attachments ...string) error
	SendPlainMessages(subject string, content string, to []string, attachments ...string) error
	SendHtmlMessage(subject string, content string, to string, attachments ...string) error
	SendHtmlMessages(subject string, content string, to []string, attachments ...string) error
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
	From     string
}

type sender struct {
	dialer *gomail.Dialer
	cfg    *Config
}

func MustNewSender(cfg *Config) Sender {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.SSL}

	return &sender{
		dialer: d,
		cfg:    cfg,
	}
}

func (s *sender) HealthCheck() bool {
	log.Debug().Msgf("check connection to smtp server %s:%d", s.cfg.Host, s.cfg.Port)

	closer, err := s.dialer.Dial()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to connect to smtp server")
		return false
	}
	if err := closer.Close(); err != nil {
		log.Error().Stack().Err(err).Msg("failed to close connection to smtp server")
		return false
	}

	return true
}

func (s *sender) SendPlainMessage(subject string, content string, to string, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	log.Debug().Msgf("sending plain email to %s", to)
	return s.dialer.DialAndSend(m)
}

func (s *sender) SendPlainMessages(subject string, content string, to []string, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	log.Debug().Msgf("sending plain email to %s", strings.Join(to, ","))
	return s.dialer.DialAndSend(m)
}

func (s *sender) SendHtmlMessage(subject string, content string, to string, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	log.Debug().Msgf("sending html email to %s", to)
	return s.dialer.DialAndSend(m)
}

func (s *sender) SendHtmlMessages(subject string, content string, to []string, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	log.Debug().Msgf("sending html email to %s", strings.Join(to, ","))
	return s.dialer.DialAndSend(m)
}
