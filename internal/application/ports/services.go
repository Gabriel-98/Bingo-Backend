package ports

import (
	"context"
	"github.com/gabriel-98/bingo-backend/internal/application/dto"
)

type AuthService interface {
	Signup(ctx context.Context, signupRequest dto.SignupRequest) (*dto.SignupResponse, error)
	Login(ctx context.Context, loginRequest dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, logoutRequest dto.LogoutRequest) error
	RefreshToken(ctx context.Context, refreshTokenRequest dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
}