package bootstrap

import (
	"fmt"

	"github.com/Ablebil/eco-sample/config"
	"github.com/Ablebil/eco-sample/internal/infra/fiber"
	"github.com/Ablebil/eco-sample/internal/infra/postgresql"
	"github.com/go-playground/validator/v10"
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

	validator := validator.New()

	app := fiber.New(cfg)
	v1 := app.Group("/api/v1")

	return app.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort))
}
