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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	log.Printf("starting")
	if err == nil {
		log.Printf("using config.json file found on disk")
	}
}

var rootCmd = &cobra.Command{
	Use:   "gl",
	Short: "go-ledger",
}

func init() {
	rootCmd.AddCommand()
}
