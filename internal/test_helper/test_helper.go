package test_helper

import (
	"os"
	"path/filepath"
)

func GetTestDataPath(filename string) string {
	// Try relative path first (for running from project root)
	if _, err := os.Stat(filepath.Join("tests/factory", filename)); err == nil {
		return filepath.Join("tests/factory", filename)
	}

	// Try from current working directory variations
	pwd, _ := os.Getwd()
	variations := []string{
		filepath.Join(pwd, "tests/factory", filename),
		filepath.Join(pwd, "../../../tests/factory", filename),
		filepath.Join(pwd, "../../../../tests/factory", filename),
		filepath.Join(pwd, "../../../../../tests/factory", filename),
	}

	for _, path := range variations {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Default to relative path (will fail if not found)
	return filepath.Join("tests/factory", filename)
}
