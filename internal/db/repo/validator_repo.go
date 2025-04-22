package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zheli/validator-key-manager-backend/pkg/models"
)

// ValidatorRepository implements the ValidatorRepo interface using SQL
type ValidatorRepository struct {
	db *sql.DB
}

// NewValidatorRepository creates a new validator repository
func NewValidatorRepository(db *sql.DB) *ValidatorRepository {
	return &ValidatorRepository{db: db}
}

// Create adds a new validator to the repository
func (r *ValidatorRepository) Create(ctx context.Context, v *models.Validator) error {
	query := `
		INSERT INTO validators (pubkey, blockchain, blockchain_network, status, client, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		v.Pubkey,
		v.Blockchain,
		v.BlockchainNetwork,
		v.Status,
		v.Client,
		now,
		now,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByPubkey retrieves a validator by its public key
func (r *ValidatorRepository) GetByPubkey(ctx context.Context, pubkey string) (*models.Validator, error) {
	query := `
		SELECT id, pubkey, blockchain, blockchain_network, status, client, created_at, updated_at
		FROM validators
		WHERE pubkey = $1`

	v := &models.Validator{}
	err := r.db.QueryRowContext(ctx, query, pubkey).Scan(
		&v.ID,
		&v.Pubkey,
		&v.Blockchain,
		&v.BlockchainNetwork,
		&v.Status,
		&v.Client,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get validator: %w", err)
	}

	return v, nil
}

// List returns a list of validators based on the provided filters
func (r *ValidatorRepository) List(ctx context.Context, filters map[string]interface{}) ([]models.Validator, error) {
	query := `
		SELECT id, pubkey, blockchain, blockchain_network, status, client, created_at, updated_at
		FROM validators
		WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add filters if provided
	if blockchain, ok := filters["blockchain"].(string); ok && blockchain != "" {
		query += fmt.Sprintf(" AND blockchain = $%d", argCount)
		args = append(args, blockchain)
		argCount++
	}

	if network, ok := filters["blockchain_network"].(string); ok && network != "" {
		query += fmt.Sprintf(" AND blockchain_network = $%d", argCount)
		args = append(args, network)
		argCount++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if client, ok := filters["client"].(string); ok && client != "" {
		query += fmt.Sprintf(" AND client = $%d", argCount)
		args = append(args, client)
		argCount++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list validators: %w", err)
	}
	defer rows.Close()

	var validators []models.Validator
	for rows.Next() {
		var v models.Validator
		err := rows.Scan(
			&v.ID,
			&v.Pubkey,
			&v.Blockchain,
			&v.BlockchainNetwork,
			&v.Status,
			&v.Client,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan validator: %w", err)
		}
		validators = append(validators, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating validators: %w", err)
	}

	return validators, nil
}

// UpdateStatus updates the status of a validator by its public key
func (r *ValidatorRepository) UpdateStatus(ctx context.Context, pubkey, status string) error {
	query := `
		UPDATE validators
		SET status = $1, updated_at = $2
		WHERE pubkey = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), pubkey)
	if err != nil {
		return fmt.Errorf("failed to update validator status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
