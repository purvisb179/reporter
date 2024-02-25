package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "reporter/docs"
	"reporter/internal/service"
)

func RegisterRoutes(r *gin.Engine, oidcService *service.OIDCService) {
	// Load HTML templates
	r.LoadHTMLGlob("internal/templates/*")

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// OIDC authentication
	authGroup := r.Group("/")
	authGroup.Use(OIDCAuthMiddleware(oidcService))
	{
		authGroup.GET("/ping", pingHandler)
	}

	// Callback handler
	r.GET("/callback", func(c *gin.Context) {
		callbackHandler(c, oidcService)
	})

	// Serve HTML route
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":         "Dropdown Example",
			"DropdownItems": []string{"Item 1", "Item 2", "Item 3"},
		})
	})
}
