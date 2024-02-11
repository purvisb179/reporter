package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all the routes for your application.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", pingHandler)
}

// pingHandler is an example handler for the /ping route.
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
