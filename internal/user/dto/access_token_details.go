package dto

// AccessTokenDetails represents the details of an access token.
type AccessTokenDetails struct {
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
}
