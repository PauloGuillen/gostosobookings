package dto

// CreatePropertyRequest represents the structure of the request body when creating a new property.
type CreatePropertyRequest struct {
	Name         string `json:"name" binding:"required"`       // Name of the property (required)
	Description  string `json:"description"`                   // Description of the property (optional)
	Address      string `json:"address"`                       // Address of the property (optional)
	ContactEmail string `json:"contact_email" binding:"email"` // Email contact (required, must be a valid email)
	ContactPhone string `json:"contact_phone"`                 // Phone contact (optional)
}
