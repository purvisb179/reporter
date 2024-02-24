package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"reporter/internal/service" // Adjust the import path according to your project structure
)

// OIDCAuthMiddleware creates a new Gin middleware for OIDC authentication
func OIDCAuthMiddleware(oidcService *service.OIDCService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Expecting the header to be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format: Bearer <token>"})
			return
		}
		rawToken := parts[1]

		// Verify the token with the OIDC service
		_, err := oidcService.VerifyToken(rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Token is valid; proceed with the request
		c.Next()
	}
}
