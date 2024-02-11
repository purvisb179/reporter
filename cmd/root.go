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
	viper.SetEnvPrefix("GL") // Prefix for environment variables
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("error reading config file: %v", err)
		}
		log.Println("starting without config.json file")
	} else {
		log.Println("using config.json file found on disk")
	}
}

var rootCmd = &cobra.Command{
	Use:   "gl",
	Short: "go-ledger",
	Long:  `go-ledger is a CLI application for managing transactions.`,
}

func init() {

}
