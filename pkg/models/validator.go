package models

import (
	"errors"
	"time"
)

// ErrNotFound is returned when a requested resource is not found
var ErrNotFound = errors.New("not found")

// Validator represents a validator in the system
type Validator struct {
	ID                int64     `json:"id" db:"id"`
	Pubkey            string    `json:"pubkey" db:"pubkey"`
	Blockchain        string    `json:"blockchain" db:"blockchain"`
	BlockchainNetwork string    `json:"blockchain_network" db:"blockchain_network"`
	Status            string    `json:"status" db:"status"`
	Client            string    `json:"client,omitempty" db:"client"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
