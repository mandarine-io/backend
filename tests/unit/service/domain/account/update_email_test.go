package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type UpdateEmailSuite struct {
	suite.Suite
}

func (s *UpdateEmailSuite) Test_Success(t provider.T) {
	t.Title("Successfully updates email")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}

	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", ctx, "email_verify", "new@example.com").Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", "email-verify", mock.Anything).Once().Return("content", nil)
	smtpSenderMock.On("SendHTMLMessage", mock.Anything, "content", mock.Anything, req.Email).Once().Return(nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().NoError(err)
	t.Require().Equal("new@example.com", resp.Email)
	t.Require().False(resp.IsEmailVerified)
}

func (s *UpdateEmailSuite) Test_UserNotFound(t provider.T) {
	t.Title("Returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(domain.ErrUserNotFound, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_FindUserByIdError(t provider.T) {
	t.Title("Returns error when FindUserById fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	dbError := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, dbError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(dbError, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_EmailNotChanged(t provider.T) {
	t.Title("Does not change email if the new email is the same as the current one")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "test@example.com",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().NoError(err)
	t.Require().Equal(req.Email, resp.Email)
	t.Require().Equal(userEntity.IsEmailVerified, resp.IsEmailVerified)
}

func (s *UpdateEmailSuite) Test_DuplicateEmail(t provider.T) {
	t.Title("Returns ErrDuplicateEmail when email already exists")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(true, nil)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(domain.ErrDuplicateEmail, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_ExistsUserByEmailError(t provider.T) {
	t.Title("Returns error when ExistsUserByEmail fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	dbError := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, dbError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(dbError, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_SetInCacheError(t provider.T) {
	t.Title("Returns error when cache set operation fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	cacheError := errors.New("cache error")

	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", ctx, "email_verify", "new@example.com").Once().Return("", cacheError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(cacheError, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_RenderContentError(t provider.T) {
	t.Title("Returns error when rendering email content fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	renderError := errors.New("render error")

	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", ctx, "email_verify", "new@example.com").Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", "email-verify", mock.Anything).Once().Return("", renderError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(renderError, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_SendHTMLMessageError(t provider.T) {
	t.Title("Returns error when sending HTML message fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	sendError := errors.New("SMTP error")

	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", ctx, "email_verify", "new@example.com").Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", "email-verify", mock.Anything).Once().Return("content", nil)
	smtpSenderMock.On("SendHTMLMessage", mock.Anything, "content", mock.Anything, req.Email).Once().Return(sendError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(domain.ErrSendEmail, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateEmailSuite) Test_UpdateUserError(t provider.T) {
	t.Title("Returns error when UpdateUser fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.UpdateEmailInput{
		Email: "new@example.com",
	}
	updateError := errors.New("database error")

	userRepoMock.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	otpServiceMock.On("GenerateAndSaveWithCode", ctx, "email_verify", "new@example.com").Once().Return("123456", nil)
	templateEngineMock.On("RenderHTML", "email-verify", mock.Anything).Once().Return("content", nil)
	smtpSenderMock.On("SendHTMLMessage", mock.Anything, "content", mock.Anything, req.Email).Once().Return(nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, updateError)

	resp, err := svc.UpdateEmail(ctx, userID, req, nil)

	t.Require().Equal(updateError, err)
	t.Require().Equal(v0.AccountOutput{}, resp)

}
