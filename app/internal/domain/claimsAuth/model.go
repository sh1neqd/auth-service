package claimsAuth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	ClientIP string    `json:"client_ip"`
	jwt.StandardClaims
}
