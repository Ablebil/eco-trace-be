package repository

import (
	"errors"

	"github.com/Ablebil/eco-sample/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChallengeRepositoryItf interface {
	GetActiveChallenges() ([]entity.Challenge, error)
	GetChallengeByID(id uuid.UUID) (*entity.Challenge, error)
	GetUserChallenges(userID uuid.UUID) ([]entity.UserChallenge, error)
	TakeChallenge(userID, challengeID uuid.UUID) error
	CompleteChallenge(userID, challengeID uuid.UUID) error
	GetUserChallenge(userID, challengeID uuid.UUID) (*entity.UserChallenge, error)
	UpdateUserExp(userID uuid.UUID, expToAdd int) error
	GetBadges() ([]entity.Badge, error)
	GetUserBadges(userID uuid.UUID) ([]entity.UserBadge, error)
	UnlockBadge(userID, badgeID uuid.UUID) error
	GetUserByID(userID uuid.UUID) (*entity.User, error)
}

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepositoryItf {
	return &ChallengeRepository{db}
}

func (r *ChallengeRepository) GetActiveChallenges() ([]entity.Challenge, error) {
	var challenges []entity.Challenge
	err := r.db.Where("is_active = ?", true).Find(&challenges).Error
	return challenges, err
}

func (r *ChallengeRepository) GetChallengeByID(id uuid.UUID) (*entity.Challenge, error) {
	var challenge entity.Challenge
	err := r.db.Where("id = ?", id).First(&challenge).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &challenge, nil
}

func (r *ChallengeRepository) GetUserChallenges(userID uuid.UUID) ([]entity.UserChallenge, error) {
	var userChallenges []entity.UserChallenge
	err := r.db.Preload("Challenge").Where("user_id = ?", userID).Find(&userChallenges).Error
	return userChallenges, err
}

func (r *ChallengeRepository) TakeChallenge(userID, challengeID uuid.UUID) error {
	userChallenge := entity.UserChallenge{
		UserID:      userID,
		ChallengeID: challengeID,
		Status:      entity.StatusOngoing,
	}
	return r.db.Create(&userChallenge).Error
}

func (r *ChallengeRepository) CompleteChallenge(userID, challengeID uuid.UUID) error {
	return r.db.Model(&entity.UserChallenge{}).
		Where("user_id = ? AND challenge_id = ?", userID, challengeID).
		Updates(map[string]interface{}{
			"status":       entity.StatusCompleted,
			"completed_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *ChallengeRepository) GetUserChallenge(userID, challengeID uuid.UUID) (*entity.UserChallenge, error) {
	var userChallenge entity.UserChallenge
	err := r.db.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&userChallenge).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &userChallenge, nil
}

func (r *ChallengeRepository) UpdateUserExp(userID uuid.UUID, expToAdd int) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("exp", gorm.Expr("exp + ?", expToAdd)).Error
}

func (r *ChallengeRepository) GetBadges() ([]entity.Badge, error) {
	var badges []entity.Badge
	err := r.db.Order("required_exp ASC").Find(&badges).Error
	return badges, err
}

func (r *ChallengeRepository) GetUserBadges(userID uuid.UUID) ([]entity.UserBadge, error) {
	var userBadges []entity.UserBadge
	err := r.db.Preload("Badge").Where("user_id = ?", userID).Find(&userBadges).Error
	return userBadges, err
}

func (r *ChallengeRepository) UnlockBadge(userID, badgeID uuid.UUID) error {
	userBadge := entity.UserBadge{
		UserID:  userID,
		BadgeID: badgeID,
	}
	return r.db.Create(&userBadge).Error
}

func (r *ChallengeRepository) GetUserByID(userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", userID).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
