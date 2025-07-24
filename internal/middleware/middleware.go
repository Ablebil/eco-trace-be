package middleware

import (
	"github.com/Ablebil/eco-sample/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareItf interface {
	Authentication(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.JWTItf
}

func NewMiddleware(jwt jwt.JWTItf) MiddlewareItf {
	return &Middleware{
		jwt: jwt,
	}
}
