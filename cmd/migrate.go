package cmd

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Handle database migrations",
	Long:  `This command allows you to manage database migrations.`,
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Long:  `Apply all up migrations to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeMigrations("up")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert all migrations",
	Long:  `Revert all migrations applied to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeMigrations("down")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(upCmd, downCmd)
}

func executeMigrations(direction string) {
	// Construct the database URL from the configuration
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")
	sslMode := "disable"

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	migrationsPath := "file://internal/db/migration"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatal("Migration initialization failed:", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while applying migrations: %v", err)
		} else {
			fmt.Println("Migrations up applied successfully")
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while reverting migrations: %v", err)
		} else {
			fmt.Println("Migrations down reverted successfully")
		}
	default:
		log.Fatalf("Invalid migration direction: %s", direction)
	}
}
