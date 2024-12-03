package app

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/handlers"
	"auth-service/app/internal/services"
	"auth-service/app/pkg/client/postgresql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func newApp() (*fiber.App, error) {
	return fiber.New(), nil
}

func StartApp(config *config.Config) {
	app, err := newApp()

	log.Println("postgresql initializing")
	db, err := postgresql.NewPostgresDB(config)
	if err != nil {
		log.Fatalf("failed to initilize db: %v", err)
	}

	log.Println("initialize services and handlers")
	service := services.NewService(db)
	handler := handlers.NewHandler(service, []byte(config.App.JwtSecret))

	log.Println("initialize routes")
	handler.InitRoutes(app)
	if err != nil {
		log.Fatalf("failed to initilize routes: %v", err)
	}
	port := config.App.Port
	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err.Error())
	}
}
