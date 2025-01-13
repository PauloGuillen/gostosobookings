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
	repo         repository.UserRepository
	jwtSecretKey string // Cached JWT secret key.
}

// NewAuthService creates a new AuthService with the necessary dependencies.
func NewAuthService(repo repository.UserRepository) *AuthService {
	// Cache the JWT secret key during initialization.
	jwtKey := config.GetEnv("JWT_SECRET_KEY", "your-secret-key")
	return &AuthService{
		repo:         repo,
		jwtSecretKey: jwtKey,
	}
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
	expirationTime := time.Now().Add(AccessTokenExpiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", user.ID),
		"role":    user.Role,
		"exp":     expirationTime,
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	// Save refresh token
	refreshTokenExpiration := time.Now().Add(RefreshTokenExpiration).Unix()
	err = s.repo.CreateRefreshToken(ctx, user.ID, refreshTokenExpiration)
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
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(s.jwtSecretKey), nil
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return dto.AccessTokenDetails{}, errors.ErrTokenExpired
			}
		}
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
		UserID:    userID,
		Role:      role,
		ExpiresAt: int64(exp), // Convert float64 to int64
	}
	return tokenDetails, nil
}

// ValidateToken validates a JWT token.
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) error {
	_, err := s.GetTokenDetails(ctx, tokenString)
	if err != nil {
		return err
	}

	return nil
}

// RevalidateToken revalidates a JWT token by issuing a new access token if the refresh token is valid.
func (s *AuthService) RevalidateToken(ctx context.Context, tokenString string) (string, error) {
	// Parse the token without validating expiration.
	parsedToken, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", errors.ErrTokenParsing
	}

	// Extract claims from the token.
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.ErrInvalidTokenClaims
	}

	// Extract and validate user_id from claims.
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.ErrInvalidTokenClaims
	}
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return "", errors.ErrInvalidTokenClaims
	}

	// Extract and validate role from claims.
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.ErrInvalidTokenClaims
	}

	// Check the refresh token's validity from the repository.
	tokenDetail, err := s.repo.FindRefreshToken(ctx, userIDInt)
	if err != nil {
		return "", err
	}
	if time.Now().Unix() > tokenDetail.ExpiresAt {
		return "", errors.ErrTokenExpired
	}

	// Generate a new access token with updated expiration.
	expirationTime := time.Now().Add(AccessTokenExpiration).Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userIDInt),
		"role":    role,
		"exp":     expirationTime,
	})
	newTokenString, err := newToken.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	return newTokenString, nil
}
