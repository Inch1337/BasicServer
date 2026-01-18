package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDb(user, password, dbname, host, port string) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	fmt.Println("Successfully connected to database")
	applyMigrations()
}

func applyMigrations() {
	sqlBytes, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		log.Fatal("Error reading migration file: ", err)
	}

	sqlQuery := string(sqlBytes)

	_, err = DB.Exec(sqlQuery)
	if err != nil {
		log.Fatal("Error applying migrations: ", err)
	}

	fmt.Println("Migrations applied successfully!")
}
