package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
)

type LogoutSuite struct {
	suite.Suite
}

func (s *LogoutSuite) Test_Success(t provider.T) {
	t.Title("Logout returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("Logout")
	t.Tags("Positive")

	ctx := context.Background()
	jti := uuid.New().String()

	jwtServiceMock.On("BanToken", ctx, jti).Once().Return(nil)

	err := svc.Logout(ctx, jti)

	t.Require().NoError(err)
}

func (s *LogoutSuite) Test_ErrBanToken(t provider.T) {
	t.Title("Logout returns BanToken error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Logout")
	t.Tags("Negative")

	ctx := context.Background()
	jti := uuid.New().String()

	expectedErr := errors.New("database error")
	jwtServiceMock.On("BanToken", ctx, jti).Once().Return(expectedErr)

	err := svc.Logout(ctx, jti)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
