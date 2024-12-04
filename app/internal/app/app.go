package app

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/handlers"
	"auth-service/app/internal/services"
	"auth-service/app/pkg/client/postgresql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func newApp() (*fiber.App, error) {
	return fiber.New(), nil
}

func StartApp(config *config.Config) {
	app, err := newApp()

	slog.Info("postgresql initializing")
	db, err := postgresql.NewPostgresDB(config)
	if err != nil {
		slog.Error("failed to initialize db: ", err.Error())
	}

	slog.Info("initialize services and handlers")
	service := services.NewService(db)
	handler := handlers.NewHandler(service, config)

	slog.Info("initialize routes")
	handler.InitRoutes(app)
	if err != nil {
		slog.Error("failed to initialize routes: ", err.Error())
	}
	port := config.App.Port
	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error(err.Error())
	}
}
