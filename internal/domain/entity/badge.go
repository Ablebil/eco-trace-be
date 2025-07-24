package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadgeType string

const (
	BadgeEcoWarrior      BadgeType = "eco_warrior"
	BadgeZeroEmission    BadgeType = "zero_emission"
	BadgeGreenHero       BadgeType = "green_hero"
	BadgeClimateChampion BadgeType = "climate_champion"
)

type Badge struct {
	ID          uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	Type        BadgeType  `gorm:"column:type;type:varchar(50);not null"`
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Description *string    `gorm:"column:description;type:text"`
	ImageURL    *string    `gorm:"column:image_url;type:text"`
	RequiredExp int        `gorm:"column:required_exp;type:int;not null"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
}

func (b *Badge) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	b.ID = id
	return
}
