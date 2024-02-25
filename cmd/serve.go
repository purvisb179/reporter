package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reporter/internal/api"
	database "reporter/internal/db"
	"reporter/internal/service"
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

	ctx := context.Background()

	providerURL := viper.GetString("oidc.provider_url")
	clientID := viper.GetString("oidc.client_id")
	clientSecret := viper.GetString("oidc.client_secret")
	redirectURL := viper.GetString("oidc.redirect_url")

	oidcService := service.NewOIDCService(ctx, providerURL, clientID, clientSecret, redirectURL)

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	api.RegisterRoutes(r, oidcService)

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
