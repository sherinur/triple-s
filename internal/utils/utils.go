package utils

import (
	"fmt"
	"os"
)

func IsValidBucketName(name string) bool {
	// TODO: Bucket names must be unique across the system.
	// TODO: Names should be between 3 and 63 characters long.
	// TODO: Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
	// TODO: Must not be formatted as an IP address (e.g., 192.168.0.1).
	// TODO: Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.

	return true
}

func CreateDir(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("error of CreateDir: dirPath is empty")
	}
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("error of CreateDir: %w", err)
	}
	return nil
}

// GetExecPath() returns path of executable file
func GetExecPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error of GetExecPath: %w", err)
	}
	return execPath, nil
}
