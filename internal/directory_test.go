package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryStructure(t *testing.T) {
	requiredDirs := []string{
		"cmd/validator-key-manager",
		"internal",
		"pkg",
		"web",
	}

	for _, dir := range requiredDirs {
		t.Run("Check "+dir+" exists", func(t *testing.T) {
			// Get the project root directory
			wd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get working directory: %v", err)
			}

			// Move up one level since we're in the internal directory
			projectRoot := filepath.Dir(wd)

			// Check if directory exists
			dirPath := filepath.Join(projectRoot, dir)
			info, err := os.Stat(dirPath)
			if err != nil {
				t.Errorf("Directory %s does not exist: %v", dir, err)
				return
			}

			if !info.IsDir() {
				t.Errorf("%s exists but is not a directory", dir)
			}
		})
	}
}
