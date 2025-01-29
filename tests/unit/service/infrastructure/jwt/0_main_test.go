package jwt

import (
	"context"
	"github.com/mandarine-io/backend/config"
	mock1 "github.com/mandarine-io/backend/internal/infrastructure/cache/mock"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/internal/service/infrastructure/jwt"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

var (
	ctx = context.Background()

	managerMock *mock1.ManagerMock
	cfg         config.JWTConfig
	svc         infrastructure.JWTService
)

func init() {
	managerMock = new(mock1.ManagerMock)
	cfg = config.JWTConfig{
		Secret:          "8O9Es3ewUadZZ0Ia+EI8IrLfNg1KpltORZdJ1q0dBjY=",
		AccessTokenTTL:  3600,
		RefreshTokenTTL: 86400,
	}
	svc = jwt.NewService(managerMock, cfg)
}

type JWTServiceSuite struct {
	suite.Suite
}

func TestJWTServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(JWTServiceSuite))
}

func (s *JWTServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(BanTokenSuite))
	s.RunSuite(t, new(GenerateTokensSuite))
	s.RunSuite(t, new(GetAccessTokenClaimsSuite))
	s.RunSuite(t, new(GetRefreshTokenClaimsSuite))
	s.RunSuite(t, new(GetTypeTokenSuite))
}
