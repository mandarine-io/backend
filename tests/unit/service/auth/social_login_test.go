package auth_test

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mandarine/internal/api/config"
	"mandarine/internal/api/persistence/model"
	mock2 "mandarine/internal/api/persistence/repo/mock"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/pkg/oauth"
	mock3 "mandarine/pkg/oauth/mock"
)

func Test_SocialLoginService_GetConsentPageUrl(t *testing.T) {
	oauthProvider := new(mock3.OAuthProviderMock)
	userRepo := new(mock2.UserRepositoryMock)
	cfg := &config.Config{Server: config.ServerConfig{ExternalOrigin: "https://example.com"}}
	service := auth.NewSocialLoginService(userRepo, oauthProvider, "provider", cfg)

	redirectUrl := "https://example.com/callback"

	t.Run("Success", func(t *testing.T) {
		oauthProvider.On("GetConsentPageUrl", redirectUrl).Return("consentUrl", "oauthState").Once()

		result := service.GetConsentPageUrl(context.Background(), redirectUrl)

		assert.Equal(t, "consentUrl", result.ConsentPageUrl)
		assert.Equal(t, "oauthState", result.OauthState)
	})
}

func Test_SocialLoginService_FetchUserInfo(t *testing.T) {
	oauthProvider := new(mock3.OAuthProviderMock)
	userRepo := new(mock2.UserRepositoryMock)
	cfg := &config.Config{Server: config.ServerConfig{ExternalOrigin: "https://example.com"}}
	service := auth.NewSocialLoginService(userRepo, oauthProvider, "provider", cfg)

	t.Run("Success", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		expectedUserInfo := oauth.UserInfo{Email: "test@example.com"}

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(&oauth2.Token{}, nil).Once()
		oauthProvider.On("GetUserInfo", mock.Anything, mock.Anything).Return(expectedUserInfo, nil).Once()

		userInfo, err := service.FetchUserInfo(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserInfo, userInfo)
	})

	t.Run("ErrorExchangingCodeToToken", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		expectedError := errors.New("exchange error")

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(nil, expectedError).Once()

		_, err := service.FetchUserInfo(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorGettingUserInfo", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		token := &oauth2.Token{}
		expectedError := errors.New("user info error")

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(token, nil).Once()
		oauthProvider.On("GetUserInfo", mock.Anything, token).Return(oauth.UserInfo{}, expectedError).Once()

		_, err := service.FetchUserInfo(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func Test_SocialLoginService_RegisterOrLogin(t *testing.T) {
	oauthProvider := new(mock3.OAuthProviderMock)
	userRepo := new(mock2.UserRepositoryMock)
	cfg := &config.Config{Security: config.SecurityConfig{JWT: config.JWTConfig{Secret: "secret"}}}
	service := auth.NewSocialLoginService(userRepo, oauthProvider, "provider", cfg)

	t.Run("Success_NewUser_UniqueUsername", func(t *testing.T) {
		userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(false, nil).Once()

		result, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("Success_NewUser_NotUniqueUsername", func(t *testing.T) {
		userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(true, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()

		result, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("Success_ExistingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com", IsEnabled: true}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(userEntity, nil).Once()

		result, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("ErrorFindingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		expectedError := errors.New("repo error")

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, expectedError).Once()

		_, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorExistingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		expectedError := errors.New("repo error")

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, expectedError).Once()

		_, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorCreatingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("create error")).Once()

		_, err := service.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
	})
}
