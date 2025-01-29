package smtp

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
}

func NewDialer(cfg Config) (*gomail.Dialer, error) {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.SSL} //nolint:gosec

	closer, err := d.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to check connection: %w", err)
	}

	if err := closer.Close(); err != nil {
		return nil, fmt.Errorf("failed to check connection : %w", err)
	}

	return d, nil
}
