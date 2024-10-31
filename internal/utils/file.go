package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func CreateFile(path string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func CreateDir(path string) error {
	exists, err := FileExists(path)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return nil
}

func RemoveDir(path string) error {
	if path == "" {
		return fmt.Errorf("filepath to the directory is empty")
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("error of RemoveDir: %w", err)
	}

	return nil
}

// DirEmpty() takes path of the directory as an argument.
// If directory empty, returns true and nil. False and nil in other case.
// Returns false and error, if error occurs.
func IsDirEmpty(path string) (bool, error) {
	isDirExits, err := FileExists(path)
	if err != nil {
		return false, err
	}

	if !isDirExits {
		return false, fmt.Errorf(path+" %w", ErrDirNotExist)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}

// GetExecPath() returns path of executable file
// Returns error, if error occurs.
func GetExecPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error of GetExecPath: %w", err)
	}
	return execPath, nil
}

// RemoveValue() takes slice of string and index as argument
// and removes the element in this index
func RemoveValue(records [][]string, index int) [][]string {
	if index < 0 || index >= len(records) {
		return records
	}

	return append(records[:index], records[index+1:]...)
}
