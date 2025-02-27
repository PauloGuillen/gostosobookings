package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Custom errors for User repository
var (
	ErrEmailAlreadyExists = errors.New("email is already in use")
	ErrInvalidData        = errors.New("invalid user data")
	ErrDatabase           = errors.New("database error")
	ErrPasswordHashing    = errors.New("error hashing password")
	ErrSonyflakeInit      = errors.New("failed to initialize Sonyflake")
	ErrSonyflakeNextID    = errors.New("error generating Sonyflake ID")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("error generating token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrTokenParsing       = errors.New("error parsing token")
	ErrTokenRequired      = errors.New("token is required")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenNotFound      = errors.New("token not found")
)

// HandleError is a centralized error handler for controllers
func HandleError(ctx *gin.Context, err error) {
	var statusCode int
	var errorMessage string

	switch err {
	case ErrEmailAlreadyExists:
		statusCode = http.StatusConflict
		errorMessage = err.Error()
	case ErrInvalidCredentials, ErrUserNotFound, ErrInvalidToken, ErrInvalidTokenClaims, ErrTokenParsing, ErrTokenRequired:
		statusCode = http.StatusUnauthorized
		errorMessage = err.Error()
	case ErrDatabase, ErrPasswordHashing, ErrSonyflakeInit, ErrSonyflakeNextID, ErrTokenGeneration:
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = "Unknown error"
	}

	ctx.JSON(statusCode, gin.H{"error": errorMessage})
}
