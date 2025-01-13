package middleware

import (
	"net/http"
	"strings"

	"github.com/PauloGuillen/gostosobookings/internal/errors"
	"github.com/PauloGuillen/gostosobookings/internal/user/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware ensures that the user is authenticated via a valid JWT token.
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate the token
		err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			if err != errors.ErrTokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
				c.Abort()
				return
			}

			// Attempt to revalidate the expired token
			newToken, revalidateErr := authService.RevalidateToken(c.Request.Context(), token)
			if revalidateErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: expired token"})
				c.Abort()
				return
			}

			// Set the new token in the Authorization header
			c.Header("Authorization", "Bearer "+newToken)
		}

		// Proceed to the next middleware or handler
		c.Next()
	}
}
