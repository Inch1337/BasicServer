package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var migrationsFS embed.FS

func InitDB(user, password, dbname, host, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Println("Successfully connected to database")
	if err := applyMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	log.Println("Migrations applied successfully")
	return db, nil
}

func applyMigrations(db *sql.DB) error {
	sqlBytes, err := migrationsFS.ReadFile("migrations/init.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}
	return nil
}
