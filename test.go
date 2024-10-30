package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CSVFile struct to work with files in CSV format
// File is a pointer to the opened file (os.File)
// Writer is a pointer to the writer
type CSVFile struct {
	File   *os.File
	Writer *csv.Writer
	Reader *csv.Reader
}

var ErrNotCSV = errors.New("the file is not in CSV format")

func OpenCSVForRead(path string) (*CSVFile, error) {
	if filepath.Ext(path) != ".csv" {
		return nil, ErrNotCSV
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can not open the file for reading: %w", err)
	}

	return &CSVFile{
		File:   file,
		Reader: csv.NewReader(file),
	}, nil
}

func OpenCSVForWrite(path string) (*CSVFile, error) {
	if filepath.Ext(path) != ".csv" {
		return nil, ErrNotCSV
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("can not open the file for writing in rewrite mode: %w", err)
	}

	return &CSVFile{
		File:   file,
		Writer: csv.NewWriter(file),
	}, nil
}

func OpenCSVForAppend(path string) (*CSVFile, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("can not open the file for writing in append mode: %w", err)
	}

	writer := csv.NewWriter(file)
	return &CSVFile{File: file, Writer: writer}, nil
}

func (csvFile *CSVFile) Close() error {
	if csvFile.Writer != nil {
		csvFile.Writer.Flush()
	}
	return csvFile.File.Close()
}

// ReadAllRecords() reads all records from the CSV-file.
// Returs two-dim slice string and error if it occurs.
func (csvFile *CSVFile) ReadAllRecords() ([][]string, error) {
	if csvFile.Reader == nil {
		return nil, errors.New("reader is not initialized")
	}
	records, err := csvFile.Reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can not read records: %w", err)
	}
	return records, nil
}

// RecordsToCSV() opens a file in overwrite mode and writes the records to it.
func (csvFile *CSVFile) RecordsToCSV(records [][]string) error {
	if csvFile.Writer == nil {
		return errors.New("writer is not initialized")
	}
	if err := csvFile.Writer.WriteAll(records); err != nil {
		return fmt.Errorf("can not write records to CSV file: %w", err)
	}
	return nil
}

// appendToCSV() appends record to the end of the CSV-file.
// Takes path of file and slice of strings as argument.
// If there is an error, returns error.
func (csvFile *CSVFile) AppendToCSV(record []string) error {
	if csvFile.Writer == nil {
		return errors.New("writer is not initialized")
	}
	if err := csvFile.Writer.Write(record); err != nil {
		return fmt.Errorf("can not append record to CSV file: %w", err)
	}
	return nil
}

// FindInCSV() searches for a record with the specified value in the given slice.
// Returns index of record and true if found, and -1, false if did not.
func FindInSlice(value string, records [][]string) (int, bool) {
	for index, record := range records {
		for _, field := range record {
			if field == value {
				return index, true
			}
		}
	}

	return -1, false
}

// RemoveValue() takes slice of string and index as argument
// and removes the element in this index
func RemoveValue(records [][]string, index int) [][]string {
	if index < 0 || index >= len(records) {
		return records
	}

	return append(records[:index], records[index+1:]...)
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
					os.Exit(1)
				}

				// file, err := OpenCSVForWrite(path)
				// if err != nil {
				// 	fmt.Println(err)
				// 	os.Exit(1)
				// }
				// file.Close()
			}
		case "check":
			if len(os.Args) > 2 {
				path := os.Args[2]
				isFileExists, err := FileExists(path)
				fmt.Println(isFileExists)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		case "parse":
			if len(os.Args) > 2 {
				path := os.Args[2]

				file, err := OpenCSVForRead(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()

				records, err := file.ReadAllRecords()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, record := range records {
					fmt.Println(record)
				}
				os.Exit(0)
			}
		case "append":
			if len(os.Args) > 2 {
				path := os.Args[2]

				file, err := OpenCSVForAppend(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()

				record := []string{"example", "2024-10-30T14:01:58+05:00", "2024-10-30T14:01:58+05:00", "active"}
				if err := file.AppendToCSV(record); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		case "find":
			if len(os.Args) > 3 {
				path := os.Args[2]
				value := os.Args[3]

				file, err := OpenCSVForRead(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()

				records, err := file.ReadAllRecords()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				i, isExists := FindInSlice(value, records)
				fmt.Println(i, isExists)

				os.Exit(0)
			}
		case "remove":
			if len(os.Args) > 3 {
				path := os.Args[2]
				value := os.Args[3]

				file, err := OpenCSVForRead(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()

				records, err := file.ReadAllRecords()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				i, isRecordExists := FindInSlice(value, records)
				if isRecordExists {
					records = RemoveValue(records, i)
				}

				file.Close()

				newFile, err := OpenCSVForWrite(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()

				if err := newFile.RecordsToCSV(records); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				os.Exit(0)
			}

		default:
			fmt.Println("Enter correct option.")
		}
	} else {
		fmt.Println("Enter correct option.")
	}
}
