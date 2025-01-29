package jwt

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	cachehelper "github.com/mandarine-io/backend/internal/util/cache"
	"github.com/rs/zerolog"
	"time"
)

const (
	jwtIssuer              = "mandarine"
	bannedTokenCachePrefix = "banned-token"
)

type svc struct {
	manager cache.Manager
	cfg     config.JWTConfig
	logger  zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(manager cache.Manager, cfg config.JWTConfig, opts ...Option) infrastructure.JWTService {
	p := &svc{
		manager: manager,
		cfg:     cfg,
		logger:  zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (s *svc) GetTypeToken(_ context.Context, token string) (string, error) {
	s.logger.Debug().Msg("get type JWT token")

	// Check token
	jwtToken, err := s.decodeAndValidateJWTToken(token)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to decode and validate token")

		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", infrastructure.ErrExpiredJWTToken
		}
		return "", infrastructure.ErrInvalidJWTToken
	}

	// Get claims
	claims, err := s.getClaimsFromJWTToken(jwtToken)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get token claims")
		return "", infrastructure.ErrInvalidJWTToken
	}

	// Get token type
	tokenType, ok := claims["type"].(string)
	if !ok {
		s.logger.Error().Stack().Err(err).Msg("failed to getting type claims")
		return "", infrastructure.ErrInvalidJWTToken
	}

	return tokenType, nil
}

func (s *svc) GetAccessTokenClaims(ctx context.Context, token string) (infrastructure.AccessTokenClaims, error) {
	s.logger.Debug().Msg("get access token claims")

	// Check token
	jwtToken, err := s.decodeAndValidateJWTToken(token)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to decode and validate token")

		if errors.Is(err, jwt.ErrTokenExpired) {
			return infrastructure.AccessTokenClaims{}, infrastructure.ErrExpiredJWTToken
		}
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get claims
	claims, err := s.getClaimsFromJWTToken(jwtToken)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get token claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get token type
	tokenType, ok := claims["type"].(string)
	if !ok {
		s.logger.Error().Stack().Err(err).Msg("failed to getting type claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	if tokenType != "access" {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("incorrect token type")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get all claims
	sub, err := claims.GetSubject()
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to getting sub claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("invalid UUID sub")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to getting exp claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting jti claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting username claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting email claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting role claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	IsPasswordTemp, ok := claims["IsPasswordTemp"].(bool)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting IsPasswordTemp claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	isEnabled, ok := claims["isEnabled"].(bool)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting isEnabled claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	isDeleted, ok := claims["isDeleted"].(bool)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting isDeleted claims")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Check if token has expired
	if exp.Unix() < time.Now().Unix() {
		s.logger.Error().Stack().Err(infrastructure.ErrExpiredJWTToken).Msg("expired jwt token")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrExpiredJWTToken
	}

	// Check if token is banned
	var bannedTokenJTI string
	err = s.manager.Get(ctx, cachehelper.CreateCacheKey(bannedTokenCachePrefix, jti), &bannedTokenJTI)

	if err != nil && !errors.Is(err, cache.ErrCacheEntryNotFound) {
		s.logger.Error().Stack().Err(err).Msg("failed to get banned token from cache")
		return infrastructure.AccessTokenClaims{}, err
	}

	if bannedTokenJTI == jti {
		s.logger.Error().Stack().Err(infrastructure.ErrBannedJWTToken).Msg("banned jwt token")
		return infrastructure.AccessTokenClaims{}, infrastructure.ErrBannedJWTToken
	}

	return infrastructure.AccessTokenClaims{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           role,
		IsPasswordTemp: IsPasswordTemp,
		IsEnabled:      isEnabled,
		IsDeleted:      isDeleted,
		JTI:            jti,
		Exp:            exp.Unix(),
	}, nil
}

func (s *svc) GetRefreshTokenClaims(ctx context.Context, token string) (infrastructure.RefreshTokenClaims, error) {
	s.logger.Debug().Msg("get refresh token claims")

	// Check token
	jwtToken, err := s.decodeAndValidateJWTToken(token)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to decode and validate token")

		if errors.Is(err, jwt.ErrTokenExpired) {
			return infrastructure.RefreshTokenClaims{}, infrastructure.ErrExpiredJWTToken
		}
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get claims
	claims, err := s.getClaimsFromJWTToken(jwtToken)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get token claims")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get token type
	tokenType, ok := claims["type"].(string)
	if !ok {
		s.logger.Error().Stack().Err(err).Msg("failed to getting type claims")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	if tokenType != "refresh" {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("incorrect token type")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Get all claims
	sub, err := claims.GetSubject()
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to getting sub claims")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("invalid UUID sub")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to getting exp claims")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		s.logger.Error().Stack().Err(infrastructure.ErrInvalidJWTToken).Msg("failed to getting jti claims")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken
	}

	// Check if token is banned
	exists, err := s.existsBannedToken(ctx, jti)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get banned token from cache")
		return infrastructure.RefreshTokenClaims{}, err
	}

	if exists {
		s.logger.Error().Stack().Err(infrastructure.ErrBannedJWTToken).Msg("banned jwt token")
		return infrastructure.RefreshTokenClaims{}, infrastructure.ErrBannedJWTToken
	}

	return infrastructure.RefreshTokenClaims{
		UserID: userID,
		JTI:    jti,
		Exp:    exp.Unix(),
	}, nil
}

func (s *svc) BanToken(ctx context.Context, jti string) error {
	s.logger.Debug().Msg("ban jwt token")

	return s.manager.SetWithExpiration(
		ctx,
		cachehelper.CreateCacheKey(bannedTokenCachePrefix, jti),
		jti,
		time.Duration(s.cfg.RefreshTokenTTL)*time.Second,
	)
}

func (s *svc) GenerateTokens(_ context.Context, userEntity *entity.User) (string, string, error) {
	s.logger.Debug().Msg("generate jwt tokens")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            jwtIssuer,
			"sub":            userEntity.ID.String(),
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Duration(s.cfg.AccessTokenTTL) * time.Second).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       userEntity.Username,
			"email":          userEntity.Email,
			"role":           userEntity.Role.Name,
			"IsPasswordTemp": userEntity.IsPasswordTemp,
			"isEnabled":      userEntity.IsEnabled,
			"isDeleted":      userEntity.DeletedAt != nil,
		},
	)
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  jwtIssuer,
			"sub":  userEntity.ID.String(),
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Duration(s.cfg.RefreshTokenTTL) * time.Second).Unix(),
			"jti":  jti,
			"type": "refresh",
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", err
	}

	refreshTokenSigned, err := refreshToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenSigned, refreshTokenSigned, nil
}

func (s *svc) decodeAndValidateJWTToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(
		token,
		s.getJWTToken(),
		jwt.WithIssuer(jwtIssuer),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims is empty")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	if sub == "" {
		return nil, errors.New("sub is empty")
	}

	iat, err := claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	if iat == nil {
		return nil, errors.New("iat is empty")
	}

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return nil, errors.New("jti is empty")
	}

	return jwtToken, nil
}

func (s *svc) getClaimsFromJWTToken(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, infrastructure.ErrInvalidJWTToken
	}

	return claims, nil
}

func (s *svc) getJWTToken() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, infrastructure.ErrInvalidJWTToken
		}
		return []byte(s.cfg.Secret), nil
	}
}

func (s *svc) existsBannedToken(ctx context.Context, jti string) (bool, error) {
	var bannedTokenJTI string
	err := s.manager.Get(ctx, cachehelper.CreateCacheKey(bannedTokenCachePrefix, jti), &bannedTokenJTI)

	if errors.Is(err, cache.ErrCacheEntryNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
