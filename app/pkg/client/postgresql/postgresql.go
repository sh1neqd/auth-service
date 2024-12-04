package postgresql

import (
	"auth-service/app/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

func NewPostgresDB(config *config.Config) (*sqlx.DB, error) {
	cfg := config.PostgresSQL
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode))
	if err != nil {
		slog.Error("failed to connect to db, err: ", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("failed to ping db, err: ", err.Error())
		return nil, err
	}

	return db, nil
}
