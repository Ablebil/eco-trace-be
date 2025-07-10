package bootstrap

import (
	"fmt"

	"github.com/Ablebil/eco-sample/config"
	"github.com/Ablebil/eco-sample/internal/infra/email"
	"github.com/Ablebil/eco-sample/internal/infra/fiber"
	"github.com/Ablebil/eco-sample/internal/infra/jwt"
	"github.com/Ablebil/eco-sample/internal/infra/oauth"
	"github.com/Ablebil/eco-sample/internal/infra/postgresql"
	"github.com/Ablebil/eco-sample/internal/infra/redis"
	"github.com/go-playground/validator/v10"

	AuthHandler "github.com/Ablebil/eco-sample/internal/app/auth/interface/rest"
	AuthUsecase "github.com/Ablebil/eco-sample/internal/app/auth/usecase"

	UserRepository "github.com/Ablebil/eco-sample/internal/app/user/repository"
)

func Start() error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := postgresql.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	), cfg)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := postgresql.Migrate(db); err != nil {
		return err
	}

	validator := validator.New()
	jwt := jwt.NewJWT(cfg)
	email := email.NewEmail(cfg)
	redis := redis.NewRedis(cfg)
	oauth := oauth.NewOAuth(cfg)

	app := fiber.New(cfg)
	v1 := app.Group("/api/v1")

	// Auth Domain
	userRepository := UserRepository.NewUserRepository(db)
	authUsecase := AuthUsecase.NewAuthUsecase(userRepository, cfg, jwt, email, redis, oauth)
	AuthHandler.NewAuthHandler(v1, validator, authUsecase, cfg)

	return app.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort))
}
