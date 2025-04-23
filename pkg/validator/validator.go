package validator

import (
	"fmt"
	"strings"
)

// ValidatePubkeyFormat checks if the provided pubkey string is valid
// It expects a hex string starting with 0x and containing 96 characters (48 bytes)
func ValidatePubkeyFormat(pubkey string) error {
	// Check if pubkey starts with 0x
	if !strings.HasPrefix(pubkey, "0x") {
		return fmt.Errorf("pubkey must start with 0x")
	}

	// Remove 0x prefix for length check
	hexPart := strings.TrimPrefix(pubkey, "0x")

	// Check length (48 bytes = 96 hex characters)
	if len(hexPart) != 96 {
		return fmt.Errorf("pubkey must be 48 bytes (96 hex characters) long, got %d characters", len(hexPart))
	}

	// Check if all characters are valid hex
	for _, c := range hexPart {
		if !isHexChar(c) {
			return fmt.Errorf("invalid hex character in pubkey: %c", c)
		}
	}

	return nil
}

// isHexChar checks if a rune is a valid hex character
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
