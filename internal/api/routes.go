package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-ledger/docs"
	"go-ledger/internal/service"
	"go-ledger/pkg/models"
	"net/http"
)

// RegisterRoutes registers all the routes for your application.
func RegisterRoutes(r *gin.Engine) {
	// Serve Swagger if available (make sure not to include it in production build)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", pingHandler)
	r.POST("/transaction", createTransaction)
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

// createTransaction handles the creation of a new transaction.
// @Summary Create a new transaction
// @Description Adds a new transaction to the database
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction"
// @Success 200 {object} models.Transaction
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transaction [post]
func createTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
