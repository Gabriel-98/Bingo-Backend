package pgrepos

import (
	"context"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
)

type RefreshTokenPgRepo struct {}

func NewRefreshTokenRepo() *RefreshTokenPgRepo {
	return &RefreshTokenPgRepo{}
}

func (repo *RefreshTokenPgRepo) Create(ctx context.Context, refreshToken entities.RefreshToken) (*entities.RefreshToken, error) {
	db, err := GetQueryExecutor(ctx)
	if err != nil {
		return nil, err
	}
	result := db.Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}

func (repo *RefreshTokenPgRepo) FindByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return nil, err
	}
	var refreshToken entities.RefreshToken
	result := db.Where("token = ?", token).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}

func (repo *RefreshTokenPgRepo) Update(ctx context.Context, token string, refreshToken entities.RefreshToken) (*entities.RefreshToken, error) {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return nil, err
	}
	refreshToken.Token = token
	result := db.Save(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}

func (repo *RefreshTokenPgRepo) Delete(ctx context.Context, token string) error {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return err
	}
	if result := db.Where("token = ?", token).Delete(&entities.RefreshToken{}); result.Error != nil {
		return result.Error
	}
	return nil
}