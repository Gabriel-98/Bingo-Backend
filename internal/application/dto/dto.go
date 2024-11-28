package dto

type SignupRequest struct {
	Username string `json:"username" example:"username"`
	Password string `json:"password" example:"MyPassword123"`
}

type SignupResponse struct {
	Id int64 `json:"id" example:"10253117"`
	Username string `json:"username" example:"jdoe65"`
}

type LoginRequest struct {
	Username string `json:"username" example:"username"`
	Password string `json:"password" example:"MyPassword123"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"51bt4584hjfh16fw5..."`
	RefreshToken string `json:"refresh_token" example:"f5t4gb61j65hf5g4d..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"f5t4gb61j65hf5g4d..."`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"f5t4gb61j65hf5g4d..."`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"51bt4584hjfh16fw5..."`
}