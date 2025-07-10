package rest

import (
	"fmt"
	"net/url"

	"github.com/Ablebil/eco-sample/config"
	"github.com/Ablebil/eco-sample/internal/app/auth/usecase"
	"github.com/Ablebil/eco-sample/internal/domain/dto"
	res "github.com/Ablebil/eco-sample/internal/infra/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	validator   *validator.Validate
	authUsecase usecase.AuthUsecaseItf
	cfg         *config.Config
}

func NewAuthHandler(authGroup fiber.Router, validator *validator.Validate, authUsecase usecase.AuthUsecaseItf, cfg *config.Config) {
	authHandler := AuthHandler{
		validator:   validator,
		authUsecase: authUsecase,
		cfg:         cfg,
	}

	authGroup = authGroup.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/verify-otp", authHandler.VerifyOTP)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Get("/google", authHandler.GoogleLogin)
	authGroup.Get("/google/callback", authHandler.GoogleCallback)
	authGroup.Post("/refresh-token", authHandler.RefreshToken)
	authGroup.Post("/logout", authHandler.Logout)
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrInternalServerError(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	if err := h.authUsecase.Register(*req); err != nil {
		return err
	}

	return res.Created(ctx, nil, res.RegisterSuccess)
}

func (h *AuthHandler) VerifyOTP(ctx *fiber.Ctx) error {
	req := new(dto.VerifyOTPRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrInternalServerError(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationsErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationsErrors)
	}

	accessToken, refreshToken, err := h.authUsecase.VerifyOTP(*req)
	if err != nil {
		return err
	}

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.VerifyOTPSuccess)
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrInternalServerError(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationsErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationsErrors)
	}

	accessToken, refreshToken, err := h.authUsecase.Login(*req)
	if err != nil {
		return err
	}

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.LoginSuccess)
}

func (h *AuthHandler) GoogleLogin(ctx *fiber.Ctx) error {
	url, err := h.authUsecase.GoogleLogin()
	if err != nil {
		return res.ErrInternalServerError(res.FailedGoogleLogin)
	}

	return ctx.Redirect(url, fiber.StatusSeeOther)
}

func (h *AuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	req := &dto.GoogleCallbackRequest{
		Code:  ctx.Query("code"),
		State: ctx.Query("state"),
		Error: ctx.Query("error"),
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationErrors)
	}

	accessToken, refreshToken, isNewUser, err := h.authUsecase.GoogleCallback(req)
	if err != nil {
		return err
	}

	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s&is_new_user=%t",
		h.cfg.FERedirectURL,
		url.QueryEscape(accessToken),
		url.QueryEscape(refreshToken),
		isNewUser)

	return ctx.Redirect(redirectURL, fiber.StatusSeeOther)
}

func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	req := new(dto.RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrInternalServerError(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationsErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationsErrors)
	}

	accessToken, refreshToken, err := h.authUsecase.RefreshToken(*req)
	if err != nil {
		return err
	}

	payload := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res.OK(ctx, payload, res.RefreshTokenSuccess)
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	req := new(dto.LogoutRequest)
	if err := ctx.BodyParser(req); err != nil {
		return res.ErrInternalServerError(res.FailedParsingRequestBody)
	}

	if err := h.validator.Struct(req); err != nil {
		validationsErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return res.ErrInternalServerError(res.FailedValidateRequest)
		}

		return res.ErrValidation(validationsErrors)
	}

	if err := h.authUsecase.Logout(*req); err != nil {
		return err
	}

	return res.OK(ctx, nil, res.LogoutSuccess)
}
