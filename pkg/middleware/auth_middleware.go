package middleware

import (
	"net/http"
	"strings"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware ensures that the user is authenticated via a valid JWT token.
// It validates the token, revalidates it if expired, and sets user details in the context for further use.
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Check the "Bearer " prefix and split the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate the token
		tokenDetail, err := authService.GetTokenDetails(c.Request.Context(), token)
		if err != nil {
			if err != errors.ErrTokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
				c.Abort()
				return
			}

			// Attempt to revalidate the expired token
			var newToken string
			newToken, tokenDetail, err = authService.RevalidateToken(c.Request.Context(), token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token expired and could not be revalidated"})
				c.Abort()
				return
			}

			// Set the new token in the Authorization header
			c.Header("Authorization", "Bearer "+newToken)
		}

		// Set the token details in the context for use in subsequent handlers
		c.Set("user_id", tokenDetail.UserID)
		c.Set("role", tokenDetail.Role)

		// Proceed to the next middleware or handler
		c.Next()
	}
}
