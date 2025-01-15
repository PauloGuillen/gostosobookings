package dto

// RefreshToken represents the data required to create a new refresh token.
type RefreshToken struct {
	ID        int64 `json:"id"`         // ID único do token de atualização
	UserID    int64 `json:"user_id"`    // ID do usuário associado
	ExpiresAt int64 `json:"expires_at"` // Data de expiração em timestamp UNIX
}
