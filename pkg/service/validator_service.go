package service

import (
	"context"

	"github.com/zheli/validator-key-manager-backend/pkg/models"
)

// ValidatorService provides business logic for validator operations
type ValidatorService struct {
	repo models.ValidatorRepo
}

// NewValidatorService creates a new validator service
func NewValidatorService(repo models.ValidatorRepo) *ValidatorService {
	return &ValidatorService{repo: repo}
}

// CreateValidator creates a new validator
func (s *ValidatorService) CreateValidator(ctx context.Context, v *models.Validator) error {
	return s.repo.Create(ctx, v)
}

// GetValidatorByPubkey retrieves a validator by its public key
func (s *ValidatorService) GetValidatorByPubkey(ctx context.Context, pubkey string) (*models.Validator, error) {
	return s.repo.GetByPubkey(ctx, pubkey)
}

// ListValidators retrieves a list of validators based on filters
func (s *ValidatorService) ListValidators(ctx context.Context, filters map[string]interface{}) ([]models.Validator, error) {
	return s.repo.List(ctx, filters)
}

// UpdateValidatorStatus updates the status of a validator
func (s *ValidatorService) UpdateValidatorStatus(ctx context.Context, pubkey, status string) error {
	return s.repo.UpdateStatus(ctx, pubkey, status)
}
