package models

import "time"

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
