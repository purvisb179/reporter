package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-ledger/internal/api"
	database "go-ledger/internal/db"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long: `Start the API server to handle requests. This command boots up the server 
and makes it listen for incoming API requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func startServer() {
	// Initialize the database connection
	database.InitDB()
	defer database.DB.Close()

	r := gin.Default()
	api.RegisterRoutes(r)

	// Load the server port from the config
	port := viper.GetString("serverPort")
	if port == "" {
		port = "8080" // Use a default port if not specified
	}
	fmt.Printf("Starting server on port %s\n", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
