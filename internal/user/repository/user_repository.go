package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/PauloGuillen/gostosobookings/config"
	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/model"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
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
	sql := `SELECT id, name, email, password_hash, role, created_at, updated_at
	FROM users WHERE email = $1`

	fmt.Println("FindByEmail sql:", sql)

	user := &model.User{}
	err := config.DB.QueryRow(ctx, sql, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrDatabase
	}

	return user, nil
}
