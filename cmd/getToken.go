package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strings"
)

// getTokenCmd represents the command to get a token
var getTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Get a bearer token using client credentials",
	Long: `This command uses the client credentials to request a bearer token 
from the OAuth2 provider configured in the application settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		getToken()
	},
}

func init() {
	rootCmd.AddCommand(getTokenCmd)
}

func getToken() {
	providerURL := viper.GetString("oidc.provider_url")
	clientID := viper.GetString("oidc.client_id")
	clientSecret := viper.GetString("oidc.client_secret")

	// Prepare the request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", providerURL+"/protocol/openid-connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Check if the access token is present
	if tokenResp.AccessToken == "" {
		fmt.Println("No access token found in the response")
		return
	}

	// Print the bearer token in the Authorization header format
	fmt.Printf("Authorization: Bearer %s\n", tokenResp.AccessToken)
}
