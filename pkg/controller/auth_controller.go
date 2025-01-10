package controller

import (
	"net/http"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
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
		errors.HandleError(ctx, err)
		return
	}

	token, err := a.authService.Login(ctx.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
