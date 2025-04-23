package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zheli/validator-key-manager-backend/pkg/mocks"
	"github.com/zheli/validator-key-manager-backend/pkg/models"
)

func TestValidatorService_CreateValidator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockValidatorRepo(ctrl)
	service := NewValidatorService(mockRepo)
	ctx := context.Background()

	validator := &models.Validator{
		Pubkey:            "0x123",
		Blockchain:        "ethereum",
		BlockchainNetwork: "mainnet",
		Status:            "active",
		Client:            "lighthouse",
	}

	mockRepo.EXPECT().Create(ctx, validator).Return(nil)

	err := service.CreateValidator(ctx, validator)
	assert.NoError(t, err)
}

func TestValidatorService_GetValidatorByPubkey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockValidatorRepo(ctrl)
	service := NewValidatorService(mockRepo)
	ctx := context.Background()

	expectedValidator := &models.Validator{
		ID:                1,
		Pubkey:            "0x123",
		Blockchain:        "ethereum",
		BlockchainNetwork: "mainnet",
		Status:            "active",
		Client:            "lighthouse",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mockRepo.EXPECT().GetByPubkey(ctx, "0x123").Return(expectedValidator, nil)

	validator, err := service.GetValidatorByPubkey(ctx, "0x123")
	assert.NoError(t, err)
	assert.Equal(t, expectedValidator, validator)
}

func TestValidatorService_ListValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockValidatorRepo(ctrl)
	service := NewValidatorService(mockRepo)
	ctx := context.Background()

	filters := map[string]interface{}{
		"blockchain": "ethereum",
		"status":     "active",
	}

	expectedValidators := []models.Validator{
		{
			ID:                1,
			Pubkey:            "0x123",
			Blockchain:        "ethereum",
			BlockchainNetwork: "mainnet",
			Status:            "active",
			Client:            "lighthouse",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	mockRepo.EXPECT().List(ctx, filters).Return(expectedValidators, nil)

	validators, err := service.ListValidators(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, expectedValidators, validators)
}

func TestValidatorService_UpdateValidatorStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockValidatorRepo(ctrl)
	service := NewValidatorService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().UpdateStatus(ctx, "0x123", "inactive").Return(nil)

	err := service.UpdateValidatorStatus(ctx, "0x123", "inactive")
	assert.NoError(t, err)
}

func TestCheckDuplicate(t *testing.T) {
	tests := []struct {
		name          string
		pubkey        string
		mockSetup     func(*gomock.Controller) *mocks.MockValidatorRepo
		expectedError error
	}{
		{
			name:   "pubkey exists",
			pubkey: "0x1234",
			mockSetup: func(ctrl *gomock.Controller) *mocks.MockValidatorRepo {
				mockRepo := mocks.NewMockValidatorRepo(ctrl)
				mockRepo.EXPECT().GetByPubkey(gomock.Any(), "0x1234").Return(&models.Validator{}, nil)
				return mockRepo
			},
			expectedError: errors.New("pubkey already exists"),
		},
		{
			name:   "pubkey does not exist",
			pubkey: "0x1234",
			mockSetup: func(ctrl *gomock.Controller) *mocks.MockValidatorRepo {
				mockRepo := mocks.NewMockValidatorRepo(ctrl)
				mockRepo.EXPECT().GetByPubkey(gomock.Any(), "0x1234").Return(nil, models.ErrNotFound)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:   "database error",
			pubkey: "0x1234",
			mockSetup: func(ctrl *gomock.Controller) *mocks.MockValidatorRepo {
				mockRepo := mocks.NewMockValidatorRepo(ctrl)
				mockRepo.EXPECT().GetByPubkey(gomock.Any(), "0x1234").Return(nil, errors.New("database error"))
				return mockRepo
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := tt.mockSetup(ctrl)
			service := NewValidatorService(mockRepo)

			err := service.CheckDuplicate(context.Background(), tt.pubkey)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
