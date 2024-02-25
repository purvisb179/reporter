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

func RegisterRoutes(r *gin.Engine, oidcService *service.OIDCService, reportService *service.ReportService) {
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
		// Serve the main page
		protected.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Title":         "Reporter",
				"DropdownItems": []string{"Item 1", "Item 2", "Item 3"},
			})
		})

		protected.GET("/content/item1", func(c *gin.Context) {
			labelValues, err := reportService.GetDistinctLabelValues()
			if err != nil {
				// Handle the error appropriately, maybe return an HTTP error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve label values"})
				return
			}
			// Pass the label values to the template
			c.HTML(http.StatusOK, "item1.html", gin.H{
				"labels": labelValues,
			})
		})

		protected.GET("/content/item2", func(c *gin.Context) {
			c.HTML(http.StatusOK, "item2.html", nil)
		})

		protected.GET("/content/item3", func(c *gin.Context) {
			c.HTML(http.StatusOK, "item3.html", nil)
		})

		protected.GET("/ping", pingHandler)

		protected.GET("/reports/download", func(c *gin.Context) {
			downloadReportHandler(c, reportService)
		})
	}
}
