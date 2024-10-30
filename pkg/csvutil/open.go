package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

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
