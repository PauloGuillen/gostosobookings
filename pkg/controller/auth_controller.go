package controller

import (
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
)

// AuthController handles user authentication
type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new AuthController instance
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login authenticates a user and generates a JWT
func (a *AuthController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authService.Login(ctx.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
