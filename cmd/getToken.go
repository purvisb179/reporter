package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
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

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Printf("Response: %s\n", body)
}
