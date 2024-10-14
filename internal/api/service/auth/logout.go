package auth

import (
	"context"
	"mandarine/internal/api/config"
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/persistence/repo"
	"time"
)

type LogoutService struct {
	banTokenRepo repo.BannedTokenRepository
	cfg          *config.Config
}

func NewLogoutService(banTokenRepo repo.BannedTokenRepository, cfg *config.Config) *LogoutService {
	return &LogoutService{
		banTokenRepo: banTokenRepo,
		cfg:          cfg,
	}
}

func (s *LogoutService) Logout(ctx context.Context, jti string) error {
	bannedToken := &model.BannedTokenEntity{
		JTI:       jti,
		ExpiredAt: time.Now().Add(time.Duration(s.cfg.Security.JWT.RefreshTokenTTL) * time.Second).Unix(),
	}
	_, err := s.banTokenRepo.CreateOrUpdateBannedToken(ctx, bannedToken)
	return err
}
