package pgrepos

import (
	"context"
	"github.com/gabriel-98/bingo-backend/internal/domain/entities"
)

type UserPgRepo struct {}

func NewUserRepo() *UserPgRepo {
	return &UserPgRepo{}
}

func (repo *UserPgRepo) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	db, err := GetQueryExecutor(ctx)
	if err != nil {
		return nil, err
	}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserPgRepo) FindById(ctx context.Context, id int64) (*entities.User, error) {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return nil, err
	}
	var user entities.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserPgRepo) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return nil, err
	}
	var user entities.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserPgRepo) Update(ctx context.Context, id int64, user entities.User) (*entities.User, error) {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return nil, err
	}
	user.Id = id
	result := db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserPgRepo) Delete(ctx context.Context, id int64) error {
	db, err := GetQueryExecutor(ctx);
	if err != nil {
		return err
	}
	if result := db.Where("id = ?", id).Delete(&entities.User{}); result.Error != nil {
		return result.Error
	}
	return nil
}