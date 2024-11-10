package model

import (
	"github.com/google/uuid"
	gormType "github.com/mandarine-io/Backend/internal/persistence/type"
	"time"
)

type MasterProfileFilter string

const (
	MasterProfileFilterDisplayName MasterProfileFilter = "display_name"
	MasterProfileFilterJob         MasterProfileFilter = "job"
	MasterProfileFilterPoint       MasterProfileFilter = "point"
	MasterProfileFilterAddress     MasterProfileFilter = "address"
)

type MasterProfileFilterPointValue struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type MasterProfileEntity struct {
	UserID      uuid.UUID      `gorm:"column:user_id;type:uuid;primaryKey;"`
	User        UserEntity     `gorm:"foreignkey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DisplayName string         `gorm:"column:display_name;type:text;not null"`
	Job         string         `gorm:"column:job;type:text;not null"`
	Description *string        `gorm:"column:description;type:text"`
	Point       gormType.Point `gorm:"column:point;type:geography(Point, 4326);not null;index:point_master_profiles_index"`
	Address     *string        `gorm:"column:address;type:text"`
	AvatarID    *string        `gorm:"column:avatar_id;type:text"`
	IsEnabled   bool           `gorm:"column:is_enabled;not null;default:true;index:is_enabled_master_profiles_index"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;type:timestamptz;default:now();autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;type:timestamptz;default:now();autoUpdateTime"`
}

func (MasterProfileEntity) TableName() string {
	return "master_profiles"
}
