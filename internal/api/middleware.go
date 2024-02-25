package api

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"reporter/internal/service"
	"reporter/pkg/models"
	"strings"
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

		// Prepare for custom audience validation
		ctx := context.Background()
		verifier := oidcService.Provider.Verifier(&oidc.Config{ClientID: "", SkipClientIDCheck: true})
		idToken, err := verifier.Verify(ctx, rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		// Decode the token claims for custom audience validation
		if err := idToken.Claims(&models.Claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
			return
		}

		clientID := oidcService.ClientId
		if !validateAudience(models.Claims, clientID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
			return
		}

		// Token is valid; proceed with the request
		c.Next()
	}
}

// validateAudience checks if the aud claim contains the clientID or if azp matches the clientID
func validateAudience(claims struct {
	Audience interface{} `json:"aud"`
	Azp      string      `json:"azp"`
}, clientID string) bool {
	switch aud := claims.Audience.(type) {
	case string:
		return aud == clientID || claims.Azp == clientID
	case []interface{}:
		for _, a := range aud {
			if str, ok := a.(string); ok && (str == clientID || claims.Azp == clientID) {
				return true
			}
		}
	}
	return false
}
