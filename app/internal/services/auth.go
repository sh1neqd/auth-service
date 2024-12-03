package services

import (
	"auth-service/app/internal/domain/authRefresh"
	"auth-service/app/internal/domain/claimsAuth"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/smtp"
	"time"
)

type AuthService struct {
	db *sqlx.DB
}

func NewAuthService(db *sqlx.DB) *AuthService {
	return &AuthService{db: db}
}

func (s AuthService) CreateRefreshToken(userID uuid.UUID) (string, error) {
	refreshToken := uuid.New().String()

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, issued_at, expires_at) 
			  VALUES (:id, :user_id, :token_hash, :issued_at, :expires_at)`

	_, err = s.db.NamedExec(query, map[string]interface{}{
		"id":         uuid.New(),
		"user_id":    userID,
		"token_hash": string(hashedToken),
		"issued_at":  time.Now(),
		"expires_at": expirationTime,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(refreshToken)), nil
}

func (s AuthService) CreateAccessToken(userID uuid.UUID, clientIP string, jwtSecret []byte) (string, error) {
	claims := &claimsAuth.Claims{
		UserID:   userID,
		ClientIP: clientIP,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func (s AuthService) GetRefreshToken(userID uuid.UUID) (*authRefresh.RefreshToken, error) {
	var refreshToken authRefresh.RefreshToken
	query := `SELECT id, user_id, token_hash, issued_at, expires_at 
			  FROM refresh_tokens 
			  WHERE user_id = $1
			  ORDER BY issued_at DESC
			  LIMIT 1`
	err := s.db.Get(&refreshToken, query, userID)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (s AuthService) SendEmailWarning(userID uuid.UUID) {

	// стоят заглушки, потенциально из другого микросервиса отправлять стоит
	userEmail := "*@example.com"

	from := ""
	password := ""

	to := []string{
		userEmail,
	}

	smtpHost := ""
	smtpPort := ""

	message := []byte("Обнаружен вход с другого ip-адреса. Если это не Вы, то обратитесь в тех поддержку")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}

	log.Println("Warning email sent successfully")
}
