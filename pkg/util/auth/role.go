package auth

const (
	RoleCustomer        = "customer"
	RoleAdmin           = "admin"
	RoleBusinessAdmin   = "business_admin"
	RoleBusinessManager = "business_manager"
)

// IsValidRole checks if the given role is valid.
func IsValidRole(role string) bool {
	switch role {
	case RoleCustomer, RoleAdmin, RoleBusinessAdmin, RoleBusinessManager:
		return true
	default:
		return false
	}
}
