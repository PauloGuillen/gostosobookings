package repository

import (
	"context"
	"database/sql"
	stdSql "database/sql"
	stdErrors "errors"

	"strings"

	"github.com/PauloGuillen/gostosobookings/config"
	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/model"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	CreateRefreshToken(ctx context.Context, userID int64, expiresAt int64) error
	DeleteRefreshToken(ctx context.Context, userID int64) error
	FindRefreshToken(ctx context.Context, userID int64) (dto.AccessTokenDetails, error)
}

// userRepository is the concrete implementation of UserRepository.
type userRepository struct{}

// NewUserRepository creates a new instance of userRepository.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Create inserts a new user into the database.
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	sql := `INSERT INTO users (id, name, email, password_hash)
	VALUES ($1, $2, $3, $4) RETURNING role, created_at, updated_at`

	err := config.DB.QueryRow(ctx, sql, user.ID, user.Name, user.Email, user.PasswordHash).
		Scan(&user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			return errors.ErrEmailAlreadyExists
		}
		return errors.ErrDatabase
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	strSql := `SELECT id, name, email, password_hash, role, created_at, updated_at
	FROM users WHERE email = $1`

	user := &model.User{}
	err := config.DB.QueryRow(ctx, strSql, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrDatabase
	}

	return user, nil
}

// CreateRefreshToken creates a new refresh token for the user.
func (r *userRepository) CreateRefreshToken(ctx context.Context, userID int64, expiresAt int64) error {
	sql := `DELETE FROM refresh_tokens where user_id = $1`
	_, err := config.DB.Exec(ctx, sql, userID)
	if err != nil {
		return errors.ErrDatabase
	}

	sql = `INSERT INTO refresh_tokens (user_id, expires_at) VALUES ($1, $2)`

	_, err = config.DB.Exec(ctx, sql, userID, expiresAt)
	if err != nil {
		return errors.ErrDatabase
	}
	return nil
}

// DeleteRefreshToken deletes the refresh token for the user.
func (r *userRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	sql := "DELETE FROM refresh_tokens WHERE user_id = $1"
	_, err := config.DB.Exec(ctx, sql, userID)
	if err != nil {
		return errors.ErrDatabase
	}

	return nil
}

// FindRefreshToken retrieves the refresh token for the user.
func (r *userRepository) FindRefreshToken(ctx context.Context, userID int64) (dto.AccessTokenDetails, error) {
	sql := `SELECT user_id, expires_at FROM refresh_tokens WHERE user_id = $1`
	tokenDetails := dto.AccessTokenDetails{}
	err := config.DB.QueryRow(ctx, sql, userID).Scan(&tokenDetails.UserID, &tokenDetails.ExpiresAt)
	if err != nil {
		if stdErrors.Is(err, stdSql.ErrNoRows) {
			return dto.AccessTokenDetails{}, errors.ErrTokenNotFound
		}
		return dto.AccessTokenDetails{}, errors.ErrDatabase
	}

	return tokenDetails, nil
}
