package model

import "time"

type BannedTokenEntity struct {
	ID        int       `gorm:"column:id;type:serial;primaryKey"`
	JTI       string    `gorm:"column:jti;type:text;unique;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null;type:timestamptz;default:now();autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;type:timestamptz;default:now();autoUpdateTime"`
	ExpiredAt int64     `gorm:"column:expired_at;not null;type:bigint;index:expired_at_banned_tokens_index"`
}

func (BannedTokenEntity) TableName() string {
	return "banned_tokens"
}
