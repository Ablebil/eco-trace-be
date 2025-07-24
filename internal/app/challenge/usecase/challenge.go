package usecase

import (
	challengeRepository "github.com/Ablebil/eco-sample/internal/app/challenge/repository"
	"github.com/Ablebil/eco-sample/internal/domain/dto"
	"github.com/Ablebil/eco-sample/internal/domain/entity"
	res "github.com/Ablebil/eco-sample/internal/infra/response"
	"github.com/google/uuid"
)

type ChallengeUsecaseItf interface {
	GetChallenges(userID uuid.UUID) ([]dto.GetChallengesResponse, *res.Err)
	TakeChallenge(userID uuid.UUID, req dto.TakeChallengeRequest) *res.Err
	CompleteChallenge(userID uuid.UUID, req dto.CompleteChallengeRequest) ([]dto.GetBadgesResponse, *res.Err)
	GetUserChallenges(userID uuid.UUID) ([]dto.GetUserChallengesResponse, *res.Err)
	GetBadges(userID uuid.UUID) ([]dto.GetBadgesResponse, *res.Err)
	GetUserStats(userID uuid.UUID) (*dto.GetUserStatsResponse, *res.Err)
}

type ChallengeUsecase struct {
	challengeRepository challengeRepository.ChallengeRepositoryItf
}

func NewChallengeUsecase(challengeRepository challengeRepository.ChallengeRepositoryItf) ChallengeUsecaseItf {
	return &ChallengeUsecase{
		challengeRepository: challengeRepository,
	}
}

func (uc *ChallengeUsecase) GetChallenges(userID uuid.UUID) ([]dto.GetChallengesResponse, *res.Err) {
	challenges, err := uc.challengeRepository.GetActiveChallenges()
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetChallenges)
	}

	var response []dto.GetChallengesResponse
	for _, challenge := range challenges {
		challengeResponse := dto.GetChallengesResponse{
			ID:          challenge.ID,
			Title:       challenge.Title,
			Description: challenge.Description,
			ExpReward:   challenge.ExpReward,
			IsActive:    challenge.IsActive,
			CreatedAt:   *challenge.CreatedAt,
		}

		userChallenge, err := uc.challengeRepository.GetUserChallenge(userID, challenge.ID)
		if err != nil {
			return nil, res.ErrInternalServerError(res.FailedGetUserChallenges)
		}

		if userChallenge != nil {
			status := string(userChallenge.Status)
			challengeResponse.Status = &status
		}

		response = append(response, challengeResponse)
	}

	return response, nil
}

func (uc *ChallengeUsecase) TakeChallenge(userID uuid.UUID, req dto.TakeChallengeRequest) *res.Err {
	challenge, err := uc.challengeRepository.GetChallengeByID(req.ChallengeID)
	if err != nil {
		return res.ErrInternalServerError(res.FailedGetChallenges)
	}

	if challenge == nil {
		return res.ErrNotFound(res.ChallengeNotFound)
	}

	if !challenge.IsActive {
		return res.ErrBadRequest(res.ChallengeNotActive)
	}

	userChallenge, err := uc.challengeRepository.GetUserChallenge(userID, req.ChallengeID)
	if err != nil {
		return res.ErrInternalServerError(res.FailedGetUserChallenges)
	}

	if userChallenge != nil {
		return res.ErrConflict(res.ChallengeAlreadyTaken)
	}

	if err := uc.challengeRepository.TakeChallenge(userID, req.ChallengeID); err != nil {
		return res.ErrInternalServerError(res.FailedTakeChallenge)
	}

	return nil
}

func (uc *ChallengeUsecase) CompleteChallenge(userID uuid.UUID, req dto.CompleteChallengeRequest) ([]dto.GetBadgesResponse, *res.Err) {
	userChallenge, err := uc.challengeRepository.GetUserChallenge(userID, req.ChallengeID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserChallenges)
	}

	if userChallenge == nil {
		return nil, res.ErrNotFound(res.ChallengeNotTaken)
	}

	if userChallenge.Status == entity.StatusCompleted {
		return nil, res.ErrConflict(res.ChallengeAlreadyCompleted)
	}

	if userChallenge.Status != entity.StatusOngoing {
		return nil, res.ErrBadRequest(res.ChallengeNotTaken)
	}

	challenge, err := uc.challengeRepository.GetChallengeByID(req.ChallengeID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetChallenges)
	}

	if challenge == nil {
		return nil, res.ErrNotFound(res.ChallengeNotFound)
	}

	if err := uc.challengeRepository.CompleteChallenge(userID, req.ChallengeID); err != nil {
		return nil, res.ErrInternalServerError(res.FailedCompleteChallenge)
	}

	if err := uc.challengeRepository.UpdateUserExp(userID, challenge.ExpReward); err != nil {
		return nil, res.ErrInternalServerError(res.FailedUpdateUserExp)
	}

	newBadges, errRes := uc.checkAndUnlockBadges(userID)
	if errRes != nil {
		return nil, errRes
	}

	return newBadges, nil
}

func (uc *ChallengeUsecase) GetUserChallenges(userID uuid.UUID) ([]dto.GetUserChallengesResponse, *res.Err) {
	userChallenges, err := uc.challengeRepository.GetUserChallenges(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserChallenges)
	}

	var response []dto.GetUserChallengesResponse
	for _, userChallenge := range userChallenges {
		challengeResponse := dto.GetUserChallengesResponse{
			ChallengeID: userChallenge.ChallengeID,
			Title:       userChallenge.Challenge.Title,
			Description: userChallenge.Challenge.Description,
			ExpReward:   userChallenge.Challenge.ExpReward,
			Status:      string(userChallenge.Status),
			CompletedAt: userChallenge.CompletedAt,
			CreatedAt:   *userChallenge.CreatedAt,
		}

		response = append(response, challengeResponse)
	}

	return response, nil
}

func (uc *ChallengeUsecase) GetBadges(userID uuid.UUID) ([]dto.GetBadgesResponse, *res.Err) {
	badges, err := uc.challengeRepository.GetBadges()
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetBadges)
	}

	userBadges, err := uc.challengeRepository.GetUserBadges(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserBadges)
	}

	unlockedBadges := make(map[uuid.UUID]*entity.UserBadge)
	for _, userBadge := range userBadges {
		unlockedBadges[userBadge.BadgeID] = &userBadge
	}

	var response []dto.GetBadgesResponse
	for _, badge := range badges {
		badgeResponse := dto.GetBadgesResponse{
			ID:          badge.ID,
			Type:        string(badge.Type),
			Name:        badge.Name,
			Description: badge.Description,
			ImageURL:    badge.ImageURL,
			RequiredExp: badge.RequiredExp,
			IsUnlocked:  false,
		}

		if userBadge, exists := unlockedBadges[badge.ID]; exists {
			badgeResponse.IsUnlocked = true
			badgeResponse.UnlockedAt = userBadge.UnlockedAt
		}

		response = append(response, badgeResponse)
	}

	return response, nil
}

func (uc *ChallengeUsecase) GetUserStats(userID uuid.UUID) (*dto.GetUserStatsResponse, *res.Err) {
	user, err := uc.challengeRepository.GetUserByID(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return nil, res.ErrNotFound(res.UserNotFound)
	}

	userChallenges, err := uc.challengeRepository.GetUserChallenges(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserChallenges)
	}

	completedCount := 0
	ongoingCount := 0
	for _, userChallenge := range userChallenges {
		switch userChallenge.Status {
		case entity.StatusCompleted:
			completedCount++
		case entity.StatusOngoing:
			ongoingCount++
		}
	}

	badges, errRes := uc.GetBadges(userID)
	if errRes != nil {
		return nil, errRes
	}

	response := &dto.GetUserStatsResponse{
		CurrentExp:      user.Exp,
		TotalChallenges: len(userChallenges),
		CompletedCount:  completedCount,
		OngoingCount:    ongoingCount,
		Badges:          badges,
	}

	return response, nil
}

func (uc *ChallengeUsecase) checkAndUnlockBadges(userID uuid.UUID) ([]dto.GetBadgesResponse, *res.Err) {
	user, err := uc.challengeRepository.GetUserByID(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedFindUser)
	}

	if user == nil {
		return nil, res.ErrNotFound(res.UserNotFound)
	}

	badges, err := uc.challengeRepository.GetBadges()
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetBadges)
	}

	userBadges, err := uc.challengeRepository.GetUserBadges(userID)
	if err != nil {
		return nil, res.ErrInternalServerError(res.FailedGetUserBadges)
	}

	unlockedBadgeIds := make(map[uuid.UUID]bool)
	for _, userBadge := range userBadges {
		unlockedBadgeIds[userBadge.BadgeID] = true
	}

	var newBadges []dto.GetBadgesResponse

	for _, badge := range badges {
		if user.Exp >= badge.RequiredExp && !unlockedBadgeIds[badge.ID] {
			if err := uc.challengeRepository.UnlockBadge(userID, badge.ID); err != nil {
				return nil, res.ErrInternalServerError(res.FailedUnlockBadge)
			}

			newBadge := dto.GetBadgesResponse{
				ID:          badge.ID,
				Type:        string(badge.Type),
				Name:        badge.Name,
				Description: badge.Description,
				ImageURL:    badge.ImageURL,
				RequiredExp: badge.RequiredExp,
				IsUnlocked:  true,
			}

			newBadges = append(newBadges, newBadge)
		}
	}

	return newBadges, nil
}
