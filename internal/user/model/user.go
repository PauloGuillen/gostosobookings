package model

import "time"

// User represents the structure of the users table in the database.
type User struct {
	ID           int64     `json:"id"` // Sonyflake ID (generated as BIGINT)
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required,email" validate:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role" default:"customer"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
