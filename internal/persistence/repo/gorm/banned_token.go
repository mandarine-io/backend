package gorm

import (
	"context"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type bannedTokenRepository struct {
	db *gorm.DB
}

func NewBannedTokenRepository(db *gorm.DB) repo.BannedTokenRepository {
	return &bannedTokenRepository{db}
}

func (b *bannedTokenRepository) CreateOrUpdateBannedToken(ctx context.Context, bannedToken *model.BannedTokenEntity) (*model.BannedTokenEntity, error) {
	log.Debug().Msg("create or update banned token")
	b.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(bannedToken)

	return bannedToken, nil
}

func (b *bannedTokenRepository) ExistsBannedTokenByJTI(ctx context.Context, jti string) (bool, error) {
	log.Debug().Msg("exists banned token by jti")
	var exists bool
	err := b.db.WithContext(ctx).
		Model(&model.BannedTokenEntity{}).
		Scopes(notExpiredTokens).
		Select("count(*) > 0").
		Where("jti = ?", jti).
		Find(&exists).
		Error
	return exists, err
}

func (b *bannedTokenRepository) DeleteExpiredBannedToken(ctx context.Context) error {
	log.Debug().Msg("delete expired banned token")
	return b.db.WithContext(ctx).Scopes(expiredTokens).Delete(&model.BannedTokenEntity{}).Error
}

func expiredTokens(db *gorm.DB) *gorm.DB {
	return db.Where("expired_at < ?", time.Now().Unix())
}

func notExpiredTokens(db *gorm.DB) *gorm.DB {
	return db.Where("expired_at >= ?", time.Now().Unix())
}
