package dto

import (
	"time"

	"github.com/google/uuid"
)

type TakeChallengeRequest struct {
	ChallengeID uuid.UUID `json:"challenge_id" validate:"required,uuid"`
}

type CompleteChallengeRequest struct {
	ChallengeID uuid.UUID `json:"challenge_id" validate:"required,uuid"`
}

type GetChallengesResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ExpReward   int       `json:"exp_reward"`
	IsActive    bool      `json:"is_active"`
	Status      *string   `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetUserChallengesResponse struct {
	ChallengeID uuid.UUID  `json:"challenge_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ExpReward   int        `json:"exp_reward"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type GetBadgesResponse struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	ImageURL    *string    `json:"image_url"`
	RequiredExp int        `json:"required_exp"`
	IsUnlocked  bool       `json:"is_unlocked"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
}

type GetUserStatsResponse struct {
	CurrentExp      int                 `json:"current_exp"`
	TotalChallenges int                 `json:"total_challenges"`
	CompletedCount  int                 `json:"completed_challenges"`
	OngoingCount    int                 `json:"ongoing_challenges"`
	Badges          []GetBadgesResponse `json:"badges"`
}
