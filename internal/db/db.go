// database/db.go

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/viper"
	"log"
)

// DB is the global database connection pool.
var DB *sql.DB

// InitDB initializes the database connection using the DATABASE_URL environment variable.
func InitDB() {
	// Construct the database connection string from components
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	// Format the connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a database connection.
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Verify the connection is alive, establishing a connection if necessary.
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}
