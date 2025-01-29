package jwt

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"time"
)

type BanTokenSuite struct {
	suite.Suite
}

func (s *BanTokenSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("JWT service")
	t.Feature("BanToken")
	t.Tags("Positive")

	jti := uuid.New().String()

	managerMock.On("SetWithExpiration", ctx, mock.Anything, jti, time.Duration(cfg.RefreshTokenTTL)*time.Second).
		Once().Return(nil)

	err := svc.BanToken(ctx, jti)

	t.Require().NoError(err)
}

func (s *BanTokenSuite) Test_ErrSettingCache(t provider.T) {
	t.Title("Returns setting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("BanToken")
	t.Tags("Negative")

	jti := uuid.New().String()

	cacheErr := errors.New("cache error")
	managerMock.On("SetWithExpiration", ctx, mock.Anything, jti, time.Duration(cfg.RefreshTokenTTL)*time.Second).
		Once().Return(cacheErr)

	err := svc.BanToken(ctx, jti)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
