package entity

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username        string     `gorm:"column:username;type:varchar(255);not null;unique"`
	Email           string     `gorm:"column:email;type:text;not null;unique"`
	Password        string     `gorm:"column:password;type:text;not null"`
	Role            RoleEntity `gorm:"foreignkey:RoleID;references:id;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RoleID          int        `gorm:"column:role_id;not null"`
	IsEnabled       bool       `gorm:"column:is_enabled;not null;default:true;index:is_enabled_users_index"`
	IsEmailVerified bool       `gorm:"column:is_email_verified;not null;default:false;index:is_email_verified_users_index"`
	IsPasswordTemp  bool       `gorm:"column:is_password_temp;not null;default:true;index:is_password_temp_index"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;type:timestamptz;default:now();autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;type:timestamptz;default:now();autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamptz;index:deleted_at_users_index"`
}

func (*User) TableName() string {
	return "users"
}

// BeforeCreate sets default role
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role.Name == RoleAdmin || u.Role.Name == RoleUser {
		return nil
	}

	var userCount int64
	if err := tx.Model(u).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		log.Debug().Msgf("set role %s for user", RoleAdmin)
		return tx.Model(&RoleEntity{}).Where("name = ?", RoleAdmin).First(&u.Role).Error
	}

	log.Debug().Msgf("set role %s for user", RoleUser)
	return tx.Model(&RoleEntity{}).Where("name = ?", RoleUser).First(&u.Role).Error
}
