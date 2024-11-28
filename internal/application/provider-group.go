package application

import (
	aports "github.com/gabriel-98/bingo-backend/internal/application/ports"
)

type ProviderGroup struct {
	passwordManager aports.PasswordManager
	authTokenManager aports.AuthTokenManager
}

func NewProviderGroup(
		passwordManager aports.PasswordManager,
		authTokenManager aports.AuthTokenManager,
		) *ProviderGroup {
			return &ProviderGroup{
				passwordManager: passwordManager,
				authTokenManager: authTokenManager,
			}
}

func (group *ProviderGroup) PasswordManager() aports.PasswordManager {
	return group.passwordManager
}

func (group *ProviderGroup) AuthTokenManager() aports.AuthTokenManager {
	return group.authTokenManager
}