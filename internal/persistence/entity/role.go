package entity

import (
	"time"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type RoleEntity struct {
	ID          int       `gorm:"column:id;type:serial;primaryKey"`
	Name        string    `gorm:"column:name;type:text;not null;unique"`
	Description *string   `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;type:timestamptz;default:now();autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;type:timestamptz;default:now();autoUpdateTime"`
}

func (*RoleEntity) TableName() string {
	return "roles"
}
