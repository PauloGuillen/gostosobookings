package dto

// AccessTokenDetails represents the details of an access token.
type AccessTokenDetails struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}
