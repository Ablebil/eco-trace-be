package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Challenge struct {
	ID          uuid.UUID  `gorm:"column:id;type:char(36);primaryKey;not null"`
	Title       string     `gorm:"column:title;type:varchar(255);not null"`
	Description *string    `gorm:"column:description;type:text"`
	ExpReward   int        `gorm:"column:exp_reward;type:int;default:0"`
	IsActive    bool       `gorm:"column:is_active;type:bool;default:true"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (c *Challenge) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	c.ID = id
	return
}
