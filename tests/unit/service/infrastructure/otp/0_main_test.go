package otp

import (
	"context"
	"github.com/mandarine-io/backend/config"
	mock1 "github.com/mandarine-io/backend/internal/infrastructure/cache/mock"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/internal/service/infrastructure/otp"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

var (
	ctx = context.Background()

	managerMock *mock1.ManagerMock
	cfg         config.OTPConfig
	svc         infrastructure.OTPService
)

func init() {
	managerMock = new(mock1.ManagerMock)
	cfg = config.OTPConfig{
		Length: 6,
		TTL:    120,
	}
	svc = otp.NewService(managerMock, cfg)
}

type OTPServiceSuite struct {
	suite.Suite
}

func TestOTPServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(OTPServiceSuite))
}

func (s *OTPServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(DeleteDataByCodeSuite))
	s.RunSuite(t, new(GenerateCodeSuite))
	s.RunSuite(t, new(GenerateAndSaveWithCodeSuite))
	s.RunSuite(t, new(GetDataByCodeSuite))
	s.RunSuite(t, new(SaveWithCodeSuite))
}
