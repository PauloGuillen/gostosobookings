package controller

import (
	"net/http"
	"strings"

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

func (a *AuthController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		errors.HandleError(ctx, errors.ErrTokenRequired)
		return
	}

	// Split by space and return the last part (token).
	parts := strings.Fields(token)
	if len(parts) > 1 {
		token = parts[1] // Returns the token part after "Bearer" or any prefix.
	}

	if err := a.authService.Logout(ctx.Request.Context(), token); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
