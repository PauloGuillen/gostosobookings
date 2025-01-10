package model

import "time"

// RefreshTokens represents a refresh token.
type RefreshTokens struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
