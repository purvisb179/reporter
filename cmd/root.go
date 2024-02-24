package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func Execute() {
	cobra.OnInitialize(initConfig)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err) // Explicitly ignoring the error
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("RE") // Prefix for environment variables
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: No config file found. Relying solely on environment variables.")
		} else {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "re",
	Short: "reporter",
	Long:  `reporter is a service for handling the creation of reports`,
}
