package application

import (
	aports "github.com/gabriel-98/bingo-backend/internal/application/ports"
)

type ServiceGroup struct {
	authService aports.AuthService
}

func NewServiceGroup(authService aports.AuthService) *ServiceGroup {
	return &ServiceGroup{
		authService: authService,
	}
}

func (group *ServiceGroup) AuthService() aports.AuthService {
	return group.authService
}