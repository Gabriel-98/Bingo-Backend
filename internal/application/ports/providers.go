package ports

import (
	"github.com/gabriel-98/bingo-backend/internal/application/types"
	"time"
)

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword string, password string) bool
}

type AuthTokenManager interface {
	NewRefreshToken(userAuthData types.UserAuthData) (string, error)
	NewAccessToken(userAuthData types.UserAuthData) (string, error)
	ValidateRefreshToken(tokenString string) (*types.UserAuthData, error)
	ValidateAccessToken(tokenString string) (*types.UserAuthData, error)
	IssuedAtAndExpiresAt(token string) (time.Time, time.Time, error)
}