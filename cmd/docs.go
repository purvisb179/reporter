package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate API documentation",
	Long:  `Automatically generates Swagger documentation for the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateDocs()
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}

func generateDocs() {
	// Execute swag command
	cmd := exec.Command("swag", "init")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Error generating Swagger documentation:", err)
		fmt.Println("Output:", string(output))
	} else {
		fmt.Println("Successfully generated Swagger documentation")
	}
}
