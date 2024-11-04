package check

import (
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/rs/zerolog/log"
)

type SmtpCheck struct {
	smtp smtp.Sender
}

func NewSmtpCheck(smtp smtp.Sender) *SmtpCheck {
	return &SmtpCheck{smtp: smtp}
}

func (r *SmtpCheck) Pass() bool {
	log.Debug().Msg("check smtp connection")
	return r.smtp.HealthCheck()
}

func (r *SmtpCheck) Name() string {
	return "smtp"
}
