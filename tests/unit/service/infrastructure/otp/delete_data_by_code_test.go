package otp

import (
	"errors"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type DeleteDataByCodeSuite struct {
	suite.Suite
}

func (s *DeleteDataByCodeSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("OTP service")
	t.Feature("DeleteDataByCode")
	t.Tags("Positive")

	prefix := "prefix"
	code := "123456"
	managerMock.On("Delete", ctx, mock.Anything).Once().Return(nil)

	err := svc.DeleteDataByCode(ctx, prefix, code)

	t.Require().NoError(err)
}

func (s *DeleteDataByCodeSuite) Test_ErrDeletingCache(t provider.T) {
	t.Title("Returns deleting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("DeleteDataByCode")
	t.Tags("Negative")

	prefix := "prefix"
	code := "123456"
	cacheErr := errors.New("cache error")

	managerMock.On("Delete", ctx, mock.Anything).Once().Return(cacheErr)

	err := svc.DeleteDataByCode(ctx, prefix, code)

	t.Require().Error(err)
}
