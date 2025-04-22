package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zheli/validator-key-manager-backend/pkg/models"
)

func TestValidatorRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewValidatorRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		validator     *models.Validator
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successful creation",
			validator: &models.Validator{
				Pubkey:            "0x123",
				Blockchain:        "ethereum",
				BlockchainNetwork: "mainnet",
				Status:            "active",
				Client:            "lighthouse",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO validators").
					WithArgs("0x123", "ethereum", "mainnet", "active", "lighthouse", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(1, time.Now(), time.Now()))
			},
			expectedError: nil,
		},
		{
			name: "duplicate pubkey",
			validator: &models.Validator{
				Pubkey:            "0x123",
				Blockchain:        "ethereum",
				BlockchainNetwork: "mainnet",
				Status:            "active",
				Client:            "lighthouse",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO validators").
					WithArgs("0x123", "ethereum", "mainnet", "active", "lighthouse", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Create(ctx, tt.validator)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.validator.ID)
				assert.NotZero(t, tt.validator.CreatedAt)
				assert.NotZero(t, tt.validator.UpdatedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidatorRepository_GetByPubkey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewValidatorRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		pubkey        string
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "successful retrieval",
			pubkey: "0x123",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pubkey", "blockchain", "blockchain_network", "status", "client", "created_at", "updated_at"}).
					AddRow(1, "0x123", "ethereum", "mainnet", "active", "lighthouse", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM validators WHERE pubkey = \\$1").
					WithArgs("0x123").
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			name:   "not found",
			pubkey: "0x123",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM validators WHERE pubkey = \\$1").
					WithArgs("0x123").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			validator, err := repo.GetByPubkey(ctx, tt.pubkey)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, validator)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validator)
				assert.Equal(t, tt.pubkey, validator.Pubkey)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidatorRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewValidatorRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		filters       map[string]interface{}
		mockSetup     func()
		expectedCount int
		expectedError error
	}{
		{
			name: "list all",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pubkey", "blockchain", "blockchain_network", "status", "client", "created_at", "updated_at"}).
					AddRow(1, "0x123", "ethereum", "mainnet", "active", "lighthouse", time.Now(), time.Now()).
					AddRow(2, "0x456", "ethereum", "mainnet", "active", "teku", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM validators WHERE 1=1").
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "filter by blockchain",
			filters: map[string]interface{}{
				"blockchain": "ethereum",
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "pubkey", "blockchain", "blockchain_network", "status", "client", "created_at", "updated_at"}).
					AddRow(1, "0x123", "ethereum", "mainnet", "active", "lighthouse", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM validators WHERE 1=1 AND blockchain = \\$1").
					WithArgs("ethereum").
					WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			validators, err := repo.List(ctx, tt.filters)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, validators)
			} else {
				assert.NoError(t, err)
				assert.Len(t, validators, tt.expectedCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidatorRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewValidatorRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		pubkey        string
		status        string
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "successful update",
			pubkey: "0x123",
			status: "inactive",
			mockSetup: func() {
				mock.ExpectExec("UPDATE validators").
					WithArgs("inactive", sqlmock.AnyArg(), "0x123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:   "not found",
			pubkey: "0x123",
			status: "inactive",
			mockSetup: func() {
				mock.ExpectExec("UPDATE validators").
					WithArgs("inactive", sqlmock.AnyArg(), "0x123").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdateStatus(ctx, tt.pubkey, tt.status)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
