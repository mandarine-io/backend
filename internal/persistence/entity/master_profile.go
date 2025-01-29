package entity

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/types"
	"gorm.io/gorm"
	"time"
)

type MasterProfile struct {
	ID          uuid.UUID   `gorm:"column:id;types:uuid;primaryKey;"`
	DisplayName string      `gorm:"column:display_name;types:text;not null"`
	Job         string      `gorm:"column:job;types:text;not null"`
	Description *string     `gorm:"column:description;types:text"`
	Point       types.Point `gorm:"column:point;types:geography(Point, 4326);not null;index:point_master_profiles_index"`
	Address     *string     `gorm:"column:address;types:text"`
	AvatarID    *string     `gorm:"column:avatar_id;types:text"`
	IsEnabled   bool        `gorm:"column:is_enabled;not null;default:true;index:is_enabled_master_profiles_index"`
	UserID      uuid.UUID   `gorm:"column:user_id;types:uuid;not null;unique"`
	User        User        `gorm:"foreignkey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time   `gorm:"column:created_at;not null;types:timestamptz;default:now();autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;not null;types:timestamptz;default:now();autoUpdateTime"`
}

func (MasterProfile) TableName() string {
	return "master_profiles"
}

type MasterProfileVector struct {
	ID                uuid.UUID     `gorm:"column:id;types:uuid;primaryKey;"`
	MasterProfileID   uuid.UUID     `gorm:"column:user_id;types:uuid;not null;unique"`
	MasterProfile     MasterProfile `gorm:"foreignkey:MasterProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DisplayNameVector string        `gorm:"column:display_name_vector;types:tsvector;not null"`
	JobVector         string        `gorm:"column:job_vector;types:tsvector;not null"`
	AddressVector     string        `gorm:"column:address_vector;types:tsvector;not null"`
}

func (*MasterProfileVector) TableName() string {
	return "master_profile_vectors"
}

func (p *MasterProfileVector) AfterCreate(tx *gorm.DB) (err error) {
	var masterProfile MasterProfile
	if err = tx.First(&masterProfile, "id = ?", p.MasterProfileID).Error; err != nil {
		return err
	}

	p.DisplayNameVector = "to_tsvector(COALESCE(display_name, ''))"
	p.JobVector = "to_tsvector(COALESCE(job, ''))"
	p.AddressVector = "to_tsvector(COALESCE(address, ''))"

	return tx.Create(p).Error
}

func (p *MasterProfileVector) AfterUpdate(tx *gorm.DB) (err error) {
	var masterProfile MasterProfile
	if err = tx.First(&masterProfile, "id = ?", p.MasterProfileID).Error; err != nil {
		return err
	}

	p.DisplayNameVector = "to_tsvector(COALESCE(display_name, ''))"
	p.JobVector = "to_tsvector(COALESCE(job, ''))"
	p.AddressVector = "to_tsvector(COALESCE(address, ''))"

	return tx.Model(&p).Updates(
		map[string]any{
			"display_name_vector": p.DisplayNameVector,
			"job_vector":          p.JobVector,
			"address_vector":      p.AddressVector,
		},
	).Error
}
