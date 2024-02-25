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
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public endpoint"})
	})

	r.GET("/callback", func(c *gin.Context) {
		callbackHandler(c, oidcService)
	})

	// Protected routes
	protected := r.Group("/")
	protected.Use(OIDCAuthMiddleware(oidcService))
	{
		protected.GET("/reports/download", func(c *gin.Context) {
			downloadReportHandler(c, reportService)
		})

		protected.GET("/ping", func(c *gin.Context) {
			pingHandler(c)
		})

		protected.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		protected.GET("/item1", func(c *gin.Context) {
			labelValues, err := reportService.GetDistinctLabelValues()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve label values"})
				return
			}
			c.HTML(http.StatusOK, "item1.html", gin.H{
				"Title":   "Item 1",
				"Content": "This is the content for Item 1.",
				"Labels":  labelValues,
			})
		})

		protected.GET("/item2", func(c *gin.Context) {
			c.HTML(http.StatusOK, "item2.html", nil)
		})
	}
}
