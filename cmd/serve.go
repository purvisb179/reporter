package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-ledger/internal/api"
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
	r := gin.Default()

	// Use the RegisterRoutes function to set up routes
	api.RegisterRoutes(r)

	// Read the server port from the configuration
	port := viper.GetString("serverPort")
	if port == "" {
		port = "8080" // Default to port 8080 if not specified
	}
	fmt.Printf("Starting server on port %s\n", port)

	// Start the server on the configured port and handle any errors
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
