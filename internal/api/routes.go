package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "reporter/docs"
	"reporter/internal/service"
)

func RegisterRoutes(r *gin.Engine, oidcService *service.OIDCService) {
	// Initialize session middleware using a cookie-based store for simplicity.
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Load HTML templates
	r.LoadHTMLGlob("internal/templates/*")

	// Swagger route for API documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	// These routes do not require authentication.
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public endpoint"})
	})

	// OIDC callback handler
	// This is a public endpoint because it's where the OIDC provider redirects after authentication.
	r.GET("/callback", func(c *gin.Context) {
		callbackHandler(c, oidcService)
	})

	// Protected routes
	// These routes require the user to be authenticated.
	protected := r.Group("/")
	protected.Use(OIDCAuthMiddleware(oidcService))
	{
		// Example of a protected route
		protected.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Title":         "Reporter",
				"DropdownItems": []string{"Item 1", "Item 2", "Item 3"},
			})
		})

		// Another example of a protected route
		protected.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
}
