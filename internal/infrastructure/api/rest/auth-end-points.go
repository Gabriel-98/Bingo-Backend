package rest

import (
	"github.com/gabriel-98/bingo-backend/internal/application/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *RestServer) SignupEndPoint(c *gin.Context) {
	// Read the request body.
	var signupRequest dto.SignupRequest
	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call to the service layer.
	authService := server.serviceGroup.AuthService()
	signupResponse, err := authService.Signup(c, signupRequest)
	if err != nil {
		// Default error for this end point (Error translation is pending)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "registration not successful",
			"message": err.Error(),
		})
		return
	}

	// Write the response body.
	c.JSON(http.StatusOK, signupResponse)
}

func (server *RestServer) LoginEndPoint(c *gin.Context) {
	// Read the request body.
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call to the service layer.
	authService := server.serviceGroup.AuthService()
	loginResponse, err := authService.Login(c, loginRequest)
	if err != nil {
		// Default error for this end point (Error translation is pending)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "unathorized",
			"message": err.Error(),
		})
		return
	}

	// Write the response body.
	c.JSON(http.StatusOK, loginResponse)
}

func (server *RestServer) LogoutEndPoint(c *gin.Context) {
	// Read the request body.
	var logoutRequest dto.LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call to the service layer.
	authService := server.serviceGroup.AuthService()
	err := authService.Logout(c, logoutRequest)
	if err != nil {
		// Default error for this end point (Error translation is pending)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "logout not successful",
			"message": err.Error(),
		})
		return
	}

	// Write the response body.
	c.JSON(http.StatusOK, gin.H{"status:": "you have logged out"})
}

func (server *RestServer) RefreshTokenEndPoint(c *gin.Context) {
	// Read the request body.
	var refreshTokenRequest dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call to the service layer.
	authService := server.serviceGroup.AuthService()
	refreshTokenResponse, err := authService.RefreshToken(c, refreshTokenRequest)
	if err != nil {
		// Default error for this end point (Error translation is pending)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "token refresh not successful",
			"message": err.Error(),
		})
		return
	}

	// Write the response body.
	c.JSON(http.StatusOK, refreshTokenResponse)
}