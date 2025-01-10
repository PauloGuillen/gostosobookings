package service

import (
	"context"
	stdErrors "errors"
	"fmt"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides authentication operations.
type AuthService struct {
	repo repository.UserRepository
}

// NewAuthService creates a new AuthService with the necessary dependencies.
func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Login authenticates a user by email and password.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Fetch user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		fmt.Println("err:", err)
		if stdErrors.Is(err, errors.ErrUserNotFound) {
			return "", errors.ErrUserNotFound
		}
		return "", errors.ErrDatabase
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	return tokenString, nil
}
