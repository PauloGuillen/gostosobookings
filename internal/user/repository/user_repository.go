package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/PauloGuillen/gostosobookings/internal/config"
	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/model"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user.
func CreateUser(userRequest dto.CreateUserRequest) (model.User, error) {
	user := model.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}

	// Initialize Sonyflake ID generator
	settings := sonyflake.Settings{}
	sf := sonyflake.NewSonyflake(settings)
	if sf == nil {
		return model.User{}, errors.ErrSonyflakeInit
	}

	// Generate unique ID
	id, err := sf.NextID()
	if err != nil {
		return model.User{}, errors.ErrSonyflakeNextID
	}
	user.ID = int64(id)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.ErrPasswordHashing
	}
	user.PasswordHash = string(hashedPassword)

	// Insert user into the database
	sql := `INSERT INTO users (id, name, email, password_hash)
	VALUES ($1, $2, $3, $4) RETURNING role, created_at, updated_at`

	err = config.DB.QueryRow(context.Background(), sql, user.ID, user.Name, user.Email, user.PasswordHash).
		Scan(&user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// If the email already exists, return a custom error
		fmt.Println("err:", err)
		fmt.Println("err.Error():", err.Error())
		if strings.Contains(err.Error(), "duplicate key value violates") {
			return model.User{}, errors.ErrEmailAlreadyExists
		}
		return model.User{}, errors.ErrDatabase
	}

	return user, nil
}
