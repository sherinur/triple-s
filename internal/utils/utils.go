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

func RemoveDir(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("error of RemoveDir: dirPath is empty")
	}

	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("error of RemoveDir: %w", err)
	}

	return nil
}

func IsDirEmpty(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("directory does not exist: %s", dirPath)
	}
	if err != nil {
		return false, fmt.Errorf("error accessing directory: %w", err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s is not a directory", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("error reading directory: %w", err)
	}

	return len(entries) == 0, nil
}

// GetExecPath() returns path of executable file
func GetExecPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error of GetExecPath: %w", err)
	}
	return execPath, nil
}
