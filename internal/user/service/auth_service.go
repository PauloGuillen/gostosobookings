package service

import (
	"context"
	stdErrors "errors"
	"fmt"
	"strconv"
	"time"

	"github.com/PauloGuillen/gostosobookings/config"
	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
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

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 2 * time.Hour
)

// Login authenticates a user by email and password.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Fetch user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
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
	jwtKey := config.GetEnv("JWT_SECRET_KEY", "your-secret-key")
	expirationTime := time.Now().Add(AccessTokenExpiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", user.ID),
		"role":    user.Role,
		"exp":     expirationTime,
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	// Save refresh token
	err = s.repo.CreateRefreshToken(ctx, user.ID, time.Now().Add(RefreshTokenExpiration))
	if err != nil {
		return "", errors.ErrDatabase
	}

	return tokenString, nil
}

// Logout invalidates a refresh token.
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	tokenDetail, err := s.GetTokenDetails(ctx, tokenString)
	if err != nil {
		return err
	}

	err = s.repo.DeleteRefreshToken(ctx, tokenDetail.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetTokenDetails extracts user_id, role, and exp from a JWT token.
func (s *AuthService) GetTokenDetails(ctx context.Context, tokenString string) (dto.AccessTokenDetails, error) {
	// Fetch the JWT secret key
	jwtKey := config.GetEnv("JWT_SECRET_KEY", "your-secret-key")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return dto.AccessTokenDetails{}, errors.ErrTokenParsing
	}

	// Validate and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return dto.AccessTokenDetails{}, errors.ErrInvalidTokenClaims
	}

	// Extract user_id
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return dto.AccessTokenDetails{}, errors.ErrInvalidTokenClaims
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.ErrInvalidTokenClaims
	}

	// Extract role
	role, ok := claims["role"].(string)
	if !ok {
		return dto.AccessTokenDetails{}, errors.ErrInvalidTokenClaims
	}

	// Extract exp (expiration)
	exp, ok := claims["exp"].(float64)
	if !ok {
		return dto.AccessTokenDetails{}, errors.ErrInvalidTokenClaims
	}

	// Return the parsed details
	tokenDetails := dto.AccessTokenDetails{
		UserID: userID,
		Role:   role,
		Exp:    int64(exp), // Convert float64 to int64
	}
	return tokenDetails, nil
}
