package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	migrationUpSQL = `CREATE TABLE IF NOT EXISTS validators (
    id SERIAL PRIMARY KEY,
    pubkey TEXT UNIQUE NOT NULL,
    blockchain TEXT NOT NULL,
    blockchain_network TEXT NOT NULL,
    status TEXT NOT NULL,
    client TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

	migrationDownSQL = `DROP TABLE IF EXISTS validators;`
)

func TestMigrations(t *testing.T) {
	ctx := context.Background()

	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "postgres"

	// Start PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get container host and port
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, port.Port(), dbName)

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err)
	defer db.Close()

	// Test connection
	err = db.Ping()
	require.NoError(t, err)

	// Run migrations
	m, err := migrate.New(
		"file://../../../migrations",
		dbURL,
	)
	require.NoError(t, err)
	defer m.Close()

	// Run up migration
	err = m.Up()
	require.NoError(t, err)

	// Verify table exists
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'validators'
		);
	`).Scan(&tableExists)
	require.NoError(t, err)
	require.True(t, tableExists)

	// Run down migration
	err = m.Down()
	require.NoError(t, err)

	// Verify table doesn't exist
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'validators'
		);
	`).Scan(&tableExists)
	require.NoError(t, err)
	require.False(t, tableExists)
}
