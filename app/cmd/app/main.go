package main

import (
	"auth-service/app/internal/app"
	"auth-service/app/internal/config"
	"context"
	"github.com/theartofdevel/logging"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()
	logger := logging.NewLogger(
		logging.WithLevel(cfg.App.LogLevel),
		logging.WithIsJSON(false),
	)
	ctx = logging.ContextWithLogger(ctx, logger)
	app.StartApp(cfg)
	logging.L(ctx).Info("application stopped")
	os.Exit(0)
}
