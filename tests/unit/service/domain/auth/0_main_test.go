package auth

import (
	"github.com/mandarine-io/backend/config"
	mock4 "github.com/mandarine-io/backend/internal/infrastructure/smtp/mock"
	mock5 "github.com/mandarine-io/backend/internal/infrastructure/template/mock"
	mock2 "github.com/mandarine-io/backend/internal/persistence/repo/mock"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/domain/auth"
	mock6 "github.com/mandarine-io/backend/internal/service/infrastructure/mock"
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/mandarine-io/backend/third_party/oauth/mock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	userRepoMock       *mock2.UserRepositoryMock
	smtpSenderMock     *mock4.SenderMock
	templateEngineMock *mock5.EngineMock
	jwtServiceMock     *mock6.JWTServiceMock
	otpServiceMock     *mock6.OTPServiceMock
	oauthProviderMock  *mock.ProviderMock
	cfg                config.Config
	svc                domain.AuthService
)

func init() {
	userRepoMock = &mock2.UserRepositoryMock{}
	smtpSenderMock = &mock4.SenderMock{}
	templateEngineMock = &mock5.EngineMock{}
	jwtServiceMock = &mock6.JWTServiceMock{}
	otpServiceMock = &mock6.OTPServiceMock{}
	cfg = config.Config{
		Security: config.SecurityConfig{
			OTP: config.OTPConfig{
				Length: 6,
				TTL:    600,
			},
		},
	}

	oauthProviderMock = &mock.ProviderMock{}
	oauthProviderMocks := make(map[string]oauth.Provider)
	oauthProviderMocks["mock"] = oauthProviderMock

	svc = auth.NewService(
		cfg,
		smtpSenderMock,
		templateEngineMock,
		userRepoMock,
		jwtServiceMock,
		otpServiceMock,
		oauthProviderMocks,
	)
}

type AuthServiceSuite struct {
	suite.Suite
}

func Test_AuthServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthServiceSuite))
}

func (s *AuthServiceSuite) Test(t provider.T) {
	t.Title("Run Auth Service tests")
	t.Epic("Auth service")

	s.RunSuite(t, new(FetchUserInfoSuite))
	s.RunSuite(t, new(GetConsentPageURLSuite))
	s.RunSuite(t, new(LoginSuite))
	s.RunSuite(t, new(LogoutSuite))
	s.RunSuite(t, new(RecoveryPasswordSuite))
	s.RunSuite(t, new(RefreshTokensSuite))
	s.RunSuite(t, new(RegisterConfirmSuite))
	s.RunSuite(t, new(RegisterOrLoginSuite))
	s.RunSuite(t, new(RegisterSuite))
	s.RunSuite(t, new(ResetPasswordSuite))
	s.RunSuite(t, new(VerifyRecoveryCodeSuite))
}
