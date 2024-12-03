package main

import (
	"auth-service/app/internal/app"
	"auth-service/app/internal/config"
)

func main() {
	app.StartApp(config.GetConfig())
}
