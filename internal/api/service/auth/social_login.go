package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/repo"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/internal/api/service/auth/mapper"
	"mandarine/pkg/logging"
	"mandarine/pkg/oauth"
	dto2 "mandarine/pkg/rest/dto"
)

var (
	ErrUserInfoNotReceived = dto2.NewI18nError("user info not received", "errors.userinfo_not_received")
)

type SocialLoginService struct {
	oauthProvider oauth.Provider
	userRepo      repo.UserRepository
	provider      string
	cfg           *config.Config
}

func NewSocialLoginService(
	userRepo repo.UserRepository, oauthProvider oauth.Provider, provider string, cfg *config.Config,
) *SocialLoginService {
	return &SocialLoginService{userRepo: userRepo, oauthProvider: oauthProvider, provider: provider, cfg: cfg}
}

//////////////////// Get consent page url ////////////////////

func (s *SocialLoginService) GetConsentPageUrl(_ context.Context, redirectUrl string) dto.GetConsentPageUrlOutput {
	slog.Info("Get consent page url")
	consentPageUrl, oauthState := s.oauthProvider.GetConsentPageUrl(redirectUrl)
	return dto.GetConsentPageUrlOutput{ConsentPageUrl: consentPageUrl, OauthState: oauthState}
}

//////////////////// Fetch user info ////////////////////

func (s *SocialLoginService) FetchUserInfo(ctx context.Context, input dto.FetchUserInfoInput) (
	oauth.UserInfo, error,
) {
	factoryErr := func(err error) (oauth.UserInfo, error) {
		slog.Error("Get user info error", logging.ErrorAttr(err))
		return oauth.UserInfo{}, err
	}

	// Exchange code to token
	slog.Info("Exchange code to token")
	socialLoginCallbackUrl := fmt.Sprintf("%s/auth/social/%s/callback/", s.cfg.Server.ExternalOrigin, s.provider)
	token, err := s.oauthProvider.ExchangeCodeToToken(ctx, input.Code, socialLoginCallbackUrl)
	if err != nil {
		return factoryErr(err)
	}

	// Get user info
	slog.Info("Get user info")
	userInfo, err := s.oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		if errors.Is(err, oauth.ErrUserInfoNotReceived) {
			return factoryErr(ErrUserInfoNotReceived)
		}
		return factoryErr(err)
	}

	return userInfo, nil
}

//////////////////// Register or login ////////////////////

func (s *SocialLoginService) RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (dto.JwtTokensOutput, error) {
	slog.Info("Register and login")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		slog.Error("Register and login error", logging.ErrorAttr(err))
		return dto.JwtTokensOutput{}, err
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, userInfo.Email, true)
	if err != nil {
		return factoryErr(err)
	}

	// Save user
	if userEntity == nil {
		slog.Info("Create new user")
		userInfo.Username, err = s.searchUniqueUsername(ctx, userInfo.Username)
		if err != nil {
			return factoryErr(err)
		}

		userEntity = mapper.MapUserInfoToUserEntity(userInfo)
		userEntity, err = s.userRepo.CreateUser(ctx, userEntity)
		if err != nil {
			return factoryErr(err)
		}
	} else if !userEntity.IsEnabled {
		return factoryErr(ErrUserIsBlocked)
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *SocialLoginService) searchUniqueUsername(ctx context.Context, defaultUsername string) (string, error) {
	slog.Info("Search unique username")
	exists, err := s.userRepo.ExistsUserByUsername(ctx, defaultUsername)
	if err != nil {
		return "", err
	}
	if !exists {
		return defaultUsername, nil
	}

	for {
		suffix, err := security.GenerateOTP(10)
		if err != nil {
			return "", err
		}

		username := fmt.Sprintf("%s_%s", defaultUsername, suffix)
		exists, err := s.userRepo.ExistsUserByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if !exists {
			return username, nil
		}

		slog.Debug("Unique username not found")
	}
}
