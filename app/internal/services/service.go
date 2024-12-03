package services

import (
	"auth-service/app/internal/domain/authRefresh"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateAccessToken(userID uuid.UUID, clientIP string, jwtSecret []byte) (string, error)
	CreateRefreshToken(userID uuid.UUID) (string, error)
	GetRefreshToken(userID uuid.UUID) (*authRefresh.RefreshToken, error)
	SendEmailWarning(userID uuid.UUID)
}

type Service struct {
	Auth
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		Auth: NewAuthService(db),
	}
}
