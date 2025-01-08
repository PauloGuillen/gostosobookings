package errors

import "errors"

// Custom errors for User repository
var (
	ErrEmailAlreadyExists = errors.New("email is already in use")
	ErrInvalidData        = errors.New("invalid user data")
	ErrDatabase           = errors.New("database error")
	ErrPasswordHashing    = errors.New("error hashing password")
	ErrSonyflakeInit      = errors.New("failed to initialize Sonyflake")
	ErrSonyflakeNextID    = errors.New("error generating Sonyflake ID")
)
