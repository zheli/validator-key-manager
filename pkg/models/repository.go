package models

import "context"

// ValidatorRepo defines the interface for validator data access
type ValidatorRepo interface {
	// Create adds a new validator to the repository
	Create(ctx context.Context, v *Validator) error

	// GetByPubkey retrieves a validator by its public key
	GetByPubkey(ctx context.Context, pubkey string) (*Validator, error)

	// List returns a list of validators based on the provided filters
	List(ctx context.Context, filters map[string]interface{}) ([]Validator, error)

	// UpdateStatus updates the status of a validator by its public key
	UpdateStatus(ctx context.Context, pubkey, status string) error
}
