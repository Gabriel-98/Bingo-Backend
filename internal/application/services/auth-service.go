package services

import (
	"context"
	"fmt"
	"github.com/gabriel-98/bingo-backend/internal/application/dto"
	aports "github.com/gabriel-98/bingo-backend/internal/application/ports"
	"github.com/gabriel-98/bingo-backend/internal/application/types"
	"github.com/gabriel-98/bingo-backend/internal/domain"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
)

type AuthService struct {
	dao *domain.DAO

	// Providers
	passwordManager aports.PasswordManager
	authTokenManager aports.AuthTokenManager
}

func NewAuthService(
		dao *domain.DAO,
		passwordManager aports.PasswordManager,
		authTokenManager aports.AuthTokenManager,
		) *AuthService {
			return &AuthService{
				dao: dao,
				passwordManager: passwordManager,
				authTokenManager: authTokenManager,
			}
}

func (s *AuthService) Signup(ctx context.Context, signupRequest dto.SignupRequest) (*dto.SignupResponse, error) {
	// Required repositories and providers.
	userRepo := s.dao.UserRepo()
	passwordManager := s.passwordManager

	// Hash password.
	hashedPassword, err := passwordManager.HashPassword(signupRequest.Password)
	if err != nil {
		return nil, err
	}

	// Create an user account and return the user id.
	user := &entities.User{
		Username: signupRequest.Username,
		Password: hashedPassword,
	}
	user, err = userRepo.Create(ctx, *user)
	if err != nil {
		return nil, err
	}
	return &dto.SignupResponse{ user.Id, user.Username } , nil
}

func (s *AuthService) Login(ctx context.Context, loginRequest dto.LoginRequest) (*dto.LoginResponse, error) {
	// Required repositories and providers.
	userRepo := s.dao.UserRepo()
	refreshTokenRepo := s.dao.RefreshTokenRepo()
	passwordManager := s.passwordManager
	authTokenManager := s.authTokenManager

	// Retrieve the user with this username.
	user, err := userRepo.FindByUsername(ctx, loginRequest.Username)
	if err != nil {
		return nil, err
	}

	// Validate password.
	if !passwordManager.CheckPassword(user.Password, loginRequest.Password) {
		return nil, fmt.Errorf("password validation failed")
	}

	// Generate refresh and access tokens.
	userAuthData := types.UserAuthData{ UserId: user.Id }
	refreshToken, err := authTokenManager.NewRefreshToken(userAuthData)
	if err != nil {
		return nil, err
	}
	accessToken, err := authTokenManager.NewAccessToken(userAuthData)
	if err != nil {
		return nil, err
	}

	// Register refresh token.
	createdAt, expiredAt, err := authTokenManager.IssuedAtAndExpiresAt(refreshToken)
	if err != nil {
		return nil, err
	}
	rt := entities.RefreshToken{
		Token: refreshToken,
		UserId: user.Id,
		CreatedAt: createdAt, 
		ExpiresAt: expiredAt,
	}
	if _, err := refreshTokenRepo.Create(ctx, rt); err !=  nil {
		return nil, err
	}

	// Build the response and return it.
	loginResponse := dto.LoginResponse{		
		RefreshToken: refreshToken,
		AccessToken: accessToken,
	}
	return &loginResponse, nil
}

func (s *AuthService) Logout(ctx context.Context, logoutRequest dto.LogoutRequest) error {
	// Required repositories and providers.
	refreshTokenRepo := s.dao.RefreshTokenRepo()

	// Validate the refresh token is registered.
	if _, err := refreshTokenRepo.FindByToken(ctx, logoutRequest.RefreshToken); err != nil {
		return fmt.Errorf("unauthenticated user: %w", err)
	}
	
	// Delete the refresh token.
	return refreshTokenRepo.Delete(ctx, logoutRequest.RefreshToken)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenRequest dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Required repositories and providers.
	refreshTokenRepo := s.dao.RefreshTokenRepo()
	authTokenManager := s.authTokenManager

	// Validate the refresh token and retrieve custom data (UserAuthData).
	// If the token is invalid (malformed, expired, etc) or not registered, authentication
	// will fail.
	refreshToken := refreshTokenRequest.RefreshToken
	userAuthData, err := authTokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("unauthenticated user: %w", err)
	}
	if _, err := refreshTokenRepo.FindByToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("unauthenticated user: %w", err)
	}

	// Generate an access token.
	accessToken, err := authTokenManager.NewAccessToken(*userAuthData)
	if err != nil {
		return nil, err
	}

	// Build the response and return it.
	refreshTokenResponse := dto.RefreshTokenResponse{		
		AccessToken: accessToken,
	}
	return &refreshTokenResponse, nil
}