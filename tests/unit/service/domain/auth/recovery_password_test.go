package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type RecoveryPasswordSuite struct {
	suite.Suite
}

func (s *RecoveryPasswordSuite) Test_Success(t provider.T) {
	t.Title("RecoveryPassword returns Success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RecoveryPassword")
	t.Tags("Positive")

	input := v0.RecoveryPasswordInput{Email: "test@example.com"}

	userRepoMock.On("ExistsUserByEmail", mock.Anything, input.Email).Return(true, nil).Once()
	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, mock.Anything, input.Email).Return("123456", nil).Once()
	smtpSenderMock.On("SendHTMLMessage", mock.Anything, mock.Anything, mock.Anything, input.Email).Return(nil).Once()
	templateEngineMock.On("RenderHTML", "recovery-password", mock.Anything).Return("email content", nil).Once()

	err := svc.RecoveryPassword(context.Background(), input, nil)

	t.Require().NoError(err)
}

func (s *RecoveryPasswordSuite) Test_UserNotFound(t provider.T) {
	t.Title("RecoveryPassword returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RecoveryPassword")
	t.Tags("Negative")

	input := v0.RecoveryPasswordInput{Email: "test@example.com"}
	userRepoMock.On("ExistsUserByEmail", mock.Anything, input.Email).Return(false, nil).Once()

	err := svc.RecoveryPassword(context.Background(), input, nil)

	t.Require().Equal(domain.ErrUserNotFound, err)
}

func (s *RecoveryPasswordSuite) Test_ErrorSettingCache(t provider.T) {
	t.Title("RecoveryPassword returns SettingCache error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RecoveryPassword")
	t.Tags("Negative")

	input := v0.RecoveryPasswordInput{Email: "test@example.com"}
	expectedErr := errors.New("cache error")

	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, mock.Anything, input.Email).Return(
		"",
		expectedErr,
	).Once()
	userRepoMock.On("ExistsUserByEmail", mock.Anything, input.Email).Return(true, nil).Once()

	err := svc.RecoveryPassword(context.Background(), input, nil)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *RecoveryPasswordSuite) Test_ErrorSendingEmail(t provider.T) {
	t.Title("RecoveryPassword returns SendingEmail error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RecoveryPassword")
	t.Tags("Negative")

	input := v0.RecoveryPasswordInput{Email: "test@example.com"}
	expectedErr := errors.New("smtp error")

	userRepoMock.On("ExistsUserByEmail", mock.Anything, input.Email).Return(true, nil).Once()
	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, mock.Anything, input.Email).Return("123456", nil).Once()
	templateEngineMock.On("RenderHTML", "recovery-password", mock.Anything).Return("email content", nil).Once()
	smtpSenderMock.On(
		"SendHTMLMessage",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		input.Email,
	).Return(expectedErr).Once()

	err := svc.RecoveryPassword(context.Background(), input, nil)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrSendEmail, err)
}
