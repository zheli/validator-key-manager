package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name          string
		databaseURL   string
		expectedURL   string
		shouldRestore bool
	}{
		{
			name:          "with database URL",
			databaseURL:   "postgres://localhost:5432/testdb?sslmode=disable",
			expectedURL:   "postgres://localhost:5432/testdb?sslmode=disable",
			shouldRestore: true,
		},
		{
			name:          "without database URL",
			databaseURL:   "",
			expectedURL:   "",
			shouldRestore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldRestore {
				oldURL := os.Getenv("DATABASE_URL")
				defer os.Setenv("DATABASE_URL", oldURL)
			}

			if tt.databaseURL != "" {
				os.Setenv("DATABASE_URL", tt.databaseURL)
			}

			cfg := NewConfig()
			assert.Equal(t, tt.expectedURL, cfg.DatabaseURL)
		})
	}
}

func TestNewDB(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "empty database URL",
			config: &Config{
				DatabaseURL: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					db.Close()
				}
			}
		})
	}
}

func TestNewDBWithMock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	// Set up expectations
	mock.ExpectPing()

	// Create a custom DB that uses our mock
	customDB := &sql.DB{}
	*customDB = *mockDB

	// Test the connection
	err = customDB.Ping()
	assert.NoError(t, err)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
