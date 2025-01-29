package smtp

import (
	"github.com/mandarine-io/backend/internal/infrastructure/smtp"
	"github.com/mandarine-io/backend/tests/integration"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	sender        smtp.Sender
	mailhogApiURL string
)

type SMTPSenderSuite struct {
	suite.Suite
}

func TestSMTPSenderSuite(t *testing.T) {
	var err error
	dialer, err := smtp.NewDialer(
		integration.Cfg.GetSMTPConfig(),
	)
	require.NoError(t, err)

	sender, err = smtp.NewSender(dialer)
	require.NoError(t, err)

	mailhogApiURL = integration.Cfg.GetMailhogAPIURL()

	suite.RunSuite(t, new(SMTPSenderSuite))
}

func (s *SMTPSenderSuite) Test(t provider.T) {
	s.RunSuite(t, new(SendPlainMessageSuite))
	s.RunSuite(t, new(SendPlainMessagesSuite))
	s.RunSuite(t, new(SendHTMLMessageSuite))
	s.RunSuite(t, new(SendHTMLMessagesSuite))
}
