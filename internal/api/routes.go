package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "reporter/docs"
	"reporter/internal/service"
)

// RegisterRoutes registers all the routes for your application.
func RegisterRoutes(r *gin.Engine, oidcService *service.OIDCService) {
	// Serve Swagger if available (make sure not to include it in production build)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply the OIDC authentication middleware to routes that require authentication
	authGroup := r.Group("/")

	authGroup.Use(OIDCAuthMiddleware(oidcService))
	{
		authGroup.GET("/ping", pingHandler)
	}

	r.GET("/callback", func(c *gin.Context) {
		callbackHandler(c, oidcService)
	})
}
