package controller

import (
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login authenticates a user and returns a JWT
func (a *AuthController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind input data
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	token, err := a.authService.Login(ctx.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return JWT
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
