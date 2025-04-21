package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	DatabaseURL string
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

// NewDB creates a new database connection
func NewDB(cfg *Config) (*sql.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
