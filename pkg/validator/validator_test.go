package validator

import (
	"testing"
)

func TestValidatePubkeyFormat(t *testing.T) {
	tests := []struct {
		name    string
		pubkey  string
		wantErr bool
	}{
		{
			name:    "valid pubkey",
			pubkey:  "0x123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
			wantErr: false,
		},
		{
			name:    "missing 0x prefix",
			pubkey:  "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
			wantErr: true,
		},
		{
			name:    "too short",
			pubkey:  "0x12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
			wantErr: true,
		},
		{
			name:    "too long",
			pubkey:  "0x1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
			wantErr: true,
		},
		{
			name:    "invalid hex character",
			pubkey:  "0x12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345g",
			wantErr: true,
		},
		{
			name:    "empty string",
			pubkey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePubkeyFormat(tt.pubkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePubkeyFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
