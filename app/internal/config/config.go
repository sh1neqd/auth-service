package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"sync"
)

type Config struct {
	PostgresSQL struct {
		Host     string `yaml:"host" env:"PSQL_HOST"  env-default:"localhost"`
		Port     int    `yaml:"port" env:"PSQL_PORT"  env-default:"5432"`
		Username string `yaml:"username" env:"PSQL_USERNAME"  env-default:"postgres"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-default:"postgres"`
		Database string `yaml:"database" env:"PSQL_DATABASE"  env-default:"postgres"`
		SSLMode  string `yaml:"sslmode" env:"SSL_MODE" env-default:"disable"`
	} `yaml:"database"`
	App struct {
		Port      int    `yaml:"port" env:"APP_PORT" env-default:"8000"`
		JwtSecret string `yaml:"jwtSecret" env:"APP_JWT_SECRET"`
		LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"app"`
	EmailSender struct {
		SenderEmail    string `yaml:"sender_email"`
		SenderPassword string `yaml:"sender_password"`
		UserEmail      string `yaml:"user_email"`
		SmtpHost       string `yaml:"smtp_host"`
		SmtpPort       string `yaml:"smtp_port"`
	} `yaml:"email_sender"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		slog.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			slog.Info(help)
			slog.Error("failed to read config", slog.String("err", err.Error()))
		}
	})
	return instance
}
