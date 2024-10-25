package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsValidBucketName(name string) bool {
	// TODO: Bucket names must be unique across the system.
	// TODO: Names should be between 3 and 63 characters long.
	// TODO: Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
	// TODO: Must not be formatted as an IP address (e.g., 192.168.0.1).
	// TODO: Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.

	return true
}

func IsUniqueBucketName(name string) bool {
	return true
}

// CreateDir() creates dir and returns error
func CreateDir(dirName string) error {
	if dirName == "" {
		return fmt.Errorf("error of CreateDir: dirName is empty")
	}

	execPath, err := GetExecPath()
	if err != nil {
		return fmt.Errorf("error of GetExecPath(): %w", err)
	}

	execDir := filepath.Dir(execPath)
	dataDirPath := filepath.Join(execDir, dirName)

	err = os.MkdirAll(dataDirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error of CreateDir: %w", err)
	}

	return nil
}

// GetExecPath() returns path of executable file
func GetExecPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return execPath, nil
}
