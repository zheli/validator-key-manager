package mocks

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zheli/validator-key-manager-backend/pkg/models"
)

func TestMockValidatorRepo(t *testing.T) {
	// This test ensures that MockValidatorRepo implements ValidatorRepo interface
	var _ models.ValidatorRepo = (*MockValidatorRepo)(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockValidatorRepo(ctrl)
	ctx := context.Background()

	// Test Create
	mock.EXPECT().Create(ctx, &models.Validator{Pubkey: "test"}).Return(nil)
	err := mock.Create(ctx, &models.Validator{Pubkey: "test"})
	assert.NoError(t, err)

	// Test GetByPubkey
	expectedValidator := &models.Validator{Pubkey: "test"}
	mock.EXPECT().GetByPubkey(ctx, "test").Return(expectedValidator, nil)
	validator, err := mock.GetByPubkey(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, expectedValidator, validator)

	// Test List
	expectedValidators := []models.Validator{{Pubkey: "test1"}, {Pubkey: "test2"}}
	mock.EXPECT().List(ctx, map[string]interface{}{"status": "active"}).Return(expectedValidators, nil)
	validators, err := mock.List(ctx, map[string]interface{}{"status": "active"})
	assert.NoError(t, err)
	assert.Equal(t, expectedValidators, validators)

	// Test UpdateStatus
	mock.EXPECT().UpdateStatus(ctx, "test", "active").Return(nil)
	err = mock.UpdateStatus(ctx, "test", "active")
	assert.NoError(t, err)
}
