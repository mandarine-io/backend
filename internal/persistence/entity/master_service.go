package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type MasterService struct {
	ID              uuid.UUID        `gorm:"column:id;types:uuid;primaryKey;"`
	Name            string           `gorm:"column:name;types:text;not null"`
	Description     *string          `gorm:"column:description;types:text"`
	MinInterval     *time.Duration   `gorm:"column:min_interval;types:interval"`
	MaxInterval     *time.Duration   `gorm:"column:max_interval;types:interval"`
	MinPrice        *decimal.Decimal `gorm:"column:min_price;types:int"`
	MaxPrice        *decimal.Decimal `gorm:"column:max_price;types:int"`
	AvatarID        *string          `gorm:"column:avatar_id;types:text"`
	MasterProfileID uuid.UUID        `gorm:"column:master_profile_id;types:uuid;not null"`
	MasterProfile   MasterProfile    `gorm:"foreignkey:MasterProfileID;references:UserID"`
	CreatedAt       time.Time        `gorm:"column:created_at;not null;types:timestamptz;default:now();autoCreateTime"`
	UpdatedAt       time.Time        `gorm:"column:updated_at;not null;types:timestamptz;default:now();autoUpdateTime"`
}

func (MasterService) TableName() string {
	return "master_services"
}
