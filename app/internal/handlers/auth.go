package handlers

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/domain/claimsAuth"
	"auth-service/app/internal/services"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	services *services.Service
	cfg      *config.Config
}

func NewHandler(services *services.Service, cfg *config.Config) *Handler {
	return &Handler{services: services, cfg: cfg}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Accept, Content-Type, Content-Length, Accept-Encoding",
	}))

	app.Get("/token", h.issueTokens)
	app.Post("/refresh", h.refreshTokens)

}

func (h *Handler) issueTokens(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	clientIP := c.IP()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user_id",
		})
	}

	accessToken, err := h.services.CreateAccessToken(userUUID, clientIP, []byte(h.cfg.App.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create access token",
		})
	}

	refreshToken, err := h.services.CreateRefreshToken(userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) refreshTokens(c *fiber.Ctx) error {
	type RefreshRequest struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var request RefreshRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	clientIP := c.IP()

	token, err := jwt.ParseWithClaims(request.AccessToken, &claimsAuth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.App.JwtSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid access token",
		})
	}

	claims, ok := token.Claims.(*claimsAuth.Claims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid access token",
		})
	}

	storedRefreshToken, err := h.services.GetRefreshToken(claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid refresh token",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve refresh token",
		})
	}

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(request.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid refresh token format",
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(storedRefreshToken.TokenHash), decodedRefreshToken) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid refresh token",
		})
	}

	if clientIP != claims.ClientIP {
		h.services.SendEmailWarning(h.cfg)
	}

	accessToken, err := h.services.CreateAccessToken(claims.UserID, clientIP, []byte(h.cfg.App.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create access token",
		})
	}

	refreshToken, err := h.services.CreateRefreshToken(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
