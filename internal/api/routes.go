package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "reporter/docs"
)

// RegisterRoutes registers all the routes for your application.
func RegisterRoutes(r *gin.Engine) {
	// Serve Swagger if available (make sure not to include it in production build)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", pingHandler)
}

// pingHandler returns a pong response
// @Summary Pong
// @Description get a pong response
// @Tags example
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
