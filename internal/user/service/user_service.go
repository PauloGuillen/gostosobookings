package service

import (
	"context"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/dto"
	"github.com/PauloGuillen/gostosobookings/internal/user/model"
	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, userRequest dto.CreateUserRequest) (model.User, error) {
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

	// Persist user using the repository
	err = s.repo.Create(ctx, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
