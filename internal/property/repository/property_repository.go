package model

import "time"

// Property represents the structure of the properties table in the database.
type Property struct {
	ID           int64     `json:"id"` // Sonyflake ID (generated as BIGINT)
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	Address      string    `json:"address"`
	ContactEmail string    `json:"contact_email" validate:"email"`
	ContactPhone string    `json:"contact_phone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
