package account

import (
	"github.com/mandarine-io/backend/config"
	mock3 "github.com/mandarine-io/backend/internal/infrastructure/smtp/mock"
	mock4 "github.com/mandarine-io/backend/internal/infrastructure/template/mock"
	mock2 "github.com/mandarine-io/backend/internal/persistence/repo/mock"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/domain/account"
	"github.com/mandarine-io/backend/internal/service/infrastructure/mock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	userRepoMock       *mock2.UserRepositoryMock
	smtpSenderMock     *mock3.SenderMock
	templateEngineMock *mock4.EngineMock
	otpServiceMock     *mock.OTPServiceMock
	cfg                config.Config
	svc                domain.AccountService
)

func init() {
	userRepoMock = &mock2.UserRepositoryMock{}
	smtpSenderMock = &mock3.SenderMock{}
	templateEngineMock = &mock4.EngineMock{}
	otpServiceMock = &mock.OTPServiceMock{}
	cfg = config.Config{
		Security: config.SecurityConfig{
			OTP: config.OTPConfig{
				Length: 6,
				TTL:    600,
			},
		},
	}
	svc = account.NewService(cfg, userRepoMock, smtpSenderMock, templateEngineMock, otpServiceMock)
}

type AccountServiceSuite struct {
	suite.Suite
}

func TestAccountServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(AccountServiceSuite))
}

func (s *AccountServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(DeleteAccountSuite))
	s.RunSuite(t, new(GetAccountSuite))
	s.RunSuite(t, new(RestoreAccountSuite))
	s.RunSuite(t, new(SetPasswordSuite))
	s.RunSuite(t, new(UpdateEmailSuite))
	s.RunSuite(t, new(UpdatePasswordSuite))
	s.RunSuite(t, new(UpdateUsernameSuite))
	s.RunSuite(t, new(VerifyEmailSuite))
}
