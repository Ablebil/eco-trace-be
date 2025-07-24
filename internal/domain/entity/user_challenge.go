package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeStatus string

const (
	StatusOngoing   ChallengeStatus = "ongoing"
	StatusCompleted ChallengeStatus = "completed"
	StatusFailed    ChallengeStatus = "failed"
)

type UserChallenge struct {
	UserID      uuid.UUID       `gorm:"column:user_id;type:char(36);primaryKey;not null"`
	ChallengeID uuid.UUID       `gorm:"column:challenge_id;type:char(36);primaryKey;not null"`
	Status      ChallengeStatus `gorm:"column:status;type:varchar(20);default:'ongoing'"`
	CompletedAt *time.Time      `gorm:"column:completed_at;type:timestamp"`
	CreatedAt   *time.Time      `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`

	User      *User      `gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`
	Challenge *Challenge `gorm:"foreignKey:challenge_id;constraint:OnDelete:CASCADE"`
}
