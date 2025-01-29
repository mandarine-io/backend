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
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type RegisterSuite struct {
	suite.Suite
}

func (s *RegisterSuite) Test_Success(t provider.T) {
	t.Title("Register returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("Register")
	t.Tags("Positive")

	req := v0.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password",
	}

	userRepoMock.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, "register", mock.Anything).Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", mock.Anything, mock.Anything).Once().Return("email content", nil)
	smtpSenderMock.On("SendHTMLMessage", mock.Anything, mock.Anything, mock.Anything, req.Email).Once().Return(nil)

	err := svc.Register(context.Background(), req, nil)

	t.Require().NoError(err)

}

func (s *RegisterSuite) Test_UserAlreadyExists(t provider.T) {
	t.Title("Register returns UserAlreadyExists error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Register")
	t.Tags("Negative")

	req := v0.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password",
	}

	userRepoMock.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(true, nil)

	err := svc.Register(context.Background(), req, nil)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrDuplicateUser, err)
}

func (s *RegisterSuite) Test_ErrorHashingPassword(t provider.T) {
	t.Title("Register returns HashingPassword error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Register")
	t.Tags("Negative")

	req := v0.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: strings.Repeat("1", 1000),
	}

	userRepoMock.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)

	err := svc.Register(context.Background(), req, nil)

	t.Require().Error(err)
	t.Require().Equal(bcrypt.ErrPasswordTooLong, err)
}

func (s *RegisterSuite) Test_ErrorSavingCache(t provider.T) {
	t.Title("Register returns SavingCache error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Register")
	t.Tags("Negative")

	req := v0.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password",
	}
	cacheErr := errors.New("cache error")

	userRepoMock.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, "register", mock.Anything).Once().Return("", cacheErr)

	err := svc.Register(context.Background(), req, nil)

	t.Require().Error(err)
	t.Require().Equal("cache error", err.Error())
}

func (s *RegisterSuite) Test_ErrorSendingEmail(t provider.T) {
	t.Title("Register returns SendingEmail error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Register")
	t.Tags("Negative")

	req := v0.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password",
	}

	userRepoMock.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", mock.Anything, "register", mock.Anything).Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", mock.Anything, mock.Anything).Once().Return("email content", nil)
	smtpSenderMock.On(
		"SendHTMLMessage",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		req.Email,
	).Once().Return(errors.New("smtp error"))

	err := svc.Register(context.Background(), req, nil)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrSendEmail, err)
}
