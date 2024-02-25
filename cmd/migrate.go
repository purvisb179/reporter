package cmd

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	database "reporter/internal/db"
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

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Long:  `Seed database with a bunch of initial data for development or testing purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		seedDatabase()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(upCmd, downCmd, seedCmd)
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

func seedDatabase() {
	// Initialize the database connection
	database.InitDB()
	defer database.DB.Close()

	// Seed the label table
	if _, err := database.DB.Exec(`INSERT INTO label (key, value) VALUES ('exampleKey1', 'exampleValue1'), ('exampleKey2', 'exampleValue2') ON CONFLICT (key, value) DO NOTHING;`); err != nil {
		log.Fatalf("Failed to seed labels: %v", err)
	}

	// Assuming we know the IDs of the labels we just inserted are 1 and 2, for simplicity
	labels := []int{1, 2} // Adjust based on actual label IDs or retrieve them from the database

	// Start a transaction for batch inserting transactions and linking them to labels
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare statement for inserting transactions
	transStmt, err := tx.Prepare("INSERT INTO transaction (id, amount) VALUES (gen_random_uuid(), $1) RETURNING id;")
	if err != nil {
		log.Fatal(err)
	}
	defer transStmt.Close()

	// Prepare statement for linking transactions to labels
	linkStmt, err := tx.Prepare("INSERT INTO transaction_label (transaction_id, label_id) VALUES ($1, $2);")
	if err != nil {
		log.Fatal(err)
	}
	defer linkStmt.Close()

	for i := 1; i <= 1000; i++ { // Adjust as needed
		var transID uuid.UUID
		err := transStmt.QueryRow(i * 100).Scan(&transID)
		if err != nil {
			log.Fatal(err)
		}

		// Attach the transaction to a label in a round-robin fashion
		labelID := labels[(i-1)%len(labels)]
		_, err = linkStmt.Exec(transID, labelID)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database seeded successfully")
}
