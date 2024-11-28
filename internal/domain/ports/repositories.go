package ports

import (
	"context"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
)

type UserRepo interface {
	Create(ctx context.Context, user entities.User) (*entities.User, error)
	FindById(ctx context.Context, id int64) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, id int64, user entities.User) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
}

type RefreshTokenRepo interface {
	Create(ctx context.Context, refreshToken entities.RefreshToken) (*entities.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	Update(ctx context.Context, token string, refreshToken entities.RefreshToken) (*entities.RefreshToken, error)
	Delete(ctx context.Context, token string) error
}