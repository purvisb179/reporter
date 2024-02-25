package api

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"reporter/internal/service"
	"reporter/pkg/models"
	"strings"
)

// OIDCAuthMiddleware creates a new Gin middleware for OIDC authentication.
func OIDCAuthMiddleware(oidcService *service.OIDCService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for a user session first.
		session := sessions.Default(c)
		if user := session.Get("user"); user != nil {
			// User is logged in, proceed with the request.
			c.Next()
			return
		}

		// For API requests, check the Authorization header for a bearer token.
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				rawToken := parts[1]

				// Prepare for custom audience validation.
				ctx := context.Background()
				verifier := oidcService.Provider.Verifier(&oidc.Config{ClientID: oidcService.ClientId, SkipClientIDCheck: true})
				idToken, err := verifier.Verify(ctx, rawToken)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
					return
				}

				// Decode the token claims for custom audience validation.
				var claims models.Claims
				if err := idToken.Claims(&claims); err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
					return
				}

				if !oidcService.ValidateAudience(claims) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
					return
				}

				// Token is valid; proceed with the request.
				c.Next()
				return
			}
		}

		// If no session or bearer token, redirect to OIDC provider for login.
		state := "yourRandomState" // This should be generated and verified securely.
		authURL := oidcService.Config.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, authURL)
		c.Abort()
	}
}
