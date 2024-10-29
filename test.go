package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotCSV = errors.New("the file is not in CSV format")
)

func ParseCSV(path string) ([][]string, error) {
	if filepath.Ext(path) != ".csv" {
		return nil, ErrNotCSV
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func WriteCSV(filepath string) error {
	return nil
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
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "make":
			if len(os.Args) > 2 {
				path := os.Args[2]
				err := CreateFile(path)
				if err != nil {
					fmt.Println(err)
				}
				os.Exit(0)
			}
		case "check":
			if len(os.Args) > 2 {
				path := os.Args[2]
				isFileExists, err := FileExists(path)
				fmt.Println(isFileExists)
				if err != nil {
					fmt.Println(err)
				}
				os.Exit(0)
			}
		case "parse":
			if len(os.Args) > 2 {
				path := os.Args[2]

				records, err := ParseCSV(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, record := range records {
					fmt.Println(record)
				}
				os.Exit(0)
			}
		}

	}
}
