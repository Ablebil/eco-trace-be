package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserBadge struct {
	UserID     uuid.UUID  `gorm:"column:user_id;type:char(36);primaryKey;not null"`
	BadgeID    uuid.UUID  `gorm:"column:badge_id;type:char(36);primaryKey;not null"`
	UnlockedAt *time.Time `gorm:"column:unlocked_at;type:timestamp;autoCreateTime"`

	User  *User  `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
	Badge *Badge `gorm:"foreignKey:badge_id;constraint:OnDelete:CASCADE"`
}
