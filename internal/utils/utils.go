package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func IsValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	validNamePattern := regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	if !validNamePattern.MatchString(name) {
		return false
	}

	ipPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	if ipPattern.MatchString(name) {
		return false
	}
	if strings.Contains(name, "--") || strings.Contains(name, "..") {
		return false
	}

	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, ".") || strings.HasSuffix(name, "-") || strings.HasSuffix(name, ".") {
		return false
	}

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
