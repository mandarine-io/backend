package smtp

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log/slog"
	"mandarine/pkg/logging"
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
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	slog.Info("Check connection to SMTP server " + addr)

	closer, err := s.dialer.Dial()
	if err != nil {
		slog.Error("SMTP client checking error", logging.ErrorAttr(err))
		return false
	}
	if err := closer.Close(); err != nil {
		slog.Error("SMTP client checking error", logging.ErrorAttr(err))
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

	slog.Debug("Sending plain email to " + to)
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

	slog.Debug("Sending plain email to " + strings.Join(to, ","))
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

	slog.Debug("Sending html email to " + to)
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

	slog.Debug("Sending html email to " + strings.Join(to, ","))
	return s.dialer.DialAndSend(m)
}
