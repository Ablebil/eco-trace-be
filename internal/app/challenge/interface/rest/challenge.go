package rest

import (
	"github.com/Ablebil/eco-sample/internal/app/challenge/usecase"
	"github.com/Ablebil/eco-sample/internal/domain/dto"
	res "github.com/Ablebil/eco-sample/internal/infra/response"
	"github.com/Ablebil/eco-sample/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChallengeHandler struct {
	validator        *validator.Validate
	challengeUsecase usecase.ChallengeUsecaseItf
}

func NewChallengeHandler(challengeGroup fiber.Router, validator *validator.Validate, challengeUsecase usecase.ChallengeUsecaseItf, middleware middleware.MiddlewareItf) {
	challengeHandler := ChallengeHandler{
		validator:        validator,
		challengeUsecase: challengeUsecase,
	}

	challengeGroup = challengeGroup.Group("/challenges")
	challengeGroup.Get("/", middleware.Authentication, challengeHandler.GetChallenges)
	challengeGroup.Post("/take", middleware.Authentication, challengeHandler.TakeChallenge)
	challengeGroup.Post("/complete", middleware.Authentication, challengeHandler.CompleteChallenge)
	challengeGroup.Get("/my", middleware.Authentication, challengeHandler.GetUserChallenges)
	challengeGroup.Get("/badges", middleware.Authentication, challengeHandler.GetBadges)
	challengeGroup.Get("/stats", middleware.Authentication, challengeHandler.GetUserStats)
}

func (h *ChallengeHandler) GetChallenges(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	challenges, errRes := h.challengeUsecase.GetChallenges(userID)
	if errRes != nil {
		return errRes
	}

	return res.OK(ctx, challenges)
}

func (h *ChallengeHandler) TakeChallenge(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	req := new(dto.TakeChallengeRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	if errRes := h.challengeUsecase.TakeChallenge(userID, *req); errRes != nil {
		return errRes
	}

	return res.OK(ctx, nil, res.TakeChallengeSuccess)
}

func (h *ChallengeHandler) CompleteChallenge(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	req := new(dto.CompleteChallengeRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrBadRequest(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	newBadges, errRes := h.challengeUsecase.CompleteChallenge(userID, *req)
	if errRes != nil {
		return errRes
	}

	payload := map[string]interface{}{
		"message":    res.CompleteChallengeSuccess,
		"new_badges": newBadges,
	}

	if len(newBadges) > 0 {
		return res.OK(ctx, payload, res.BadgeUnlockedSuccess)
	}

	return res.OK(ctx, payload, res.CompleteChallengeSuccess)
}

func (h *ChallengeHandler) GetUserChallenges(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	challenges, errRes := h.challengeUsecase.GetUserChallenges(userID)
	if errRes != nil {
		return errRes
	}

	return res.OK(ctx, challenges)
}

func (h *ChallengeHandler) GetBadges(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	badges, errRes := h.challengeUsecase.GetBadges(userID)
	if errRes != nil {
		return errRes
	}

	return res.OK(ctx, badges)
}

func (h *ChallengeHandler) GetUserStats(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	stats, errRes := h.challengeUsecase.GetUserStats(userID)
	if errRes != nil {
		return errRes
	}

	return res.OK(ctx, stats)
}

func getUserIDFromContext(ctx *fiber.Ctx) (uuid.UUID, *res.Err) {
	userIDStr := ctx.Locals("user_id")
	if userIDStr == nil {
		return uuid.Nil, res.ErrUnauthorized("User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, res.ErrUnauthorized("Invalid user ID")
	}

	return userID, nil
}
