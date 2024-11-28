package domain

import (
	domports "github.com/gabriel-98/bingo-backend/internal/domain/ports"
)

type DAO struct {
	userRepo domports.UserRepo
	refreshTokenRepo domports.RefreshTokenRepo
}

func NewDAO(
		userRepo domports.UserRepo,
		refreshTokenRepo domports.RefreshTokenRepo,
		) *DAO {
		return &DAO{
			userRepo: userRepo,
			refreshTokenRepo: refreshTokenRepo,
		}
}

func (dao *DAO) UserRepo() domports.UserRepo {
	return dao.userRepo
}

func (dao *DAO) RefreshTokenRepo() domports.RefreshTokenRepo {
	return dao.refreshTokenRepo
}