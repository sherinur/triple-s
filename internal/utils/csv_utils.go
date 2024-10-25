package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"triple-s/internal/types"
)

func CreateBucketMeta(name string) error {
	bucket := types.NewBucket(name)

	csvWriter, file, err := openCSV("buckets.csv")
	if err != nil {
		return fmt.Errorf("error of opening or creating a bucket metadata: %w", err)
	}
	defer file.Close()
	defer csvWriter.Flush()

	data := convertBucketToArr(bucket)
	if err := csvWriter.Write(data); err != nil {
		return fmt.Errorf("error of writing a bucket metadata: %w", err)
	}

	return nil
}

func FindBucketByName(name string, records [][]string) bool {
	for _, line := range records {
		if line[0] == name {
			return true
		}
	}

	return false
}

func ParseCSV(filename string) ([][]string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(filepath.Dir(execPath), "data", filename)

	file, err := os.Open(filepath)
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

func openCSV(name string) (*csv.Writer, *os.File, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, nil, fmt.Errorf("error of getting exec path: %w", err)
	}
	dataDirPath := filepath.Join(filepath.Dir(execPath), "data")

	if err := os.MkdirAll(dataDirPath, os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("error of dir creation 'data': %w", err)
	}

	filePath := filepath.Join(dataDirPath, name)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("error of opening or creating file - 'buckets.csv': %w", err)
	}

	bufferedWriter := bufio.NewWriter(file)

	csvWriter := csv.NewWriter(bufferedWriter)

	return csvWriter, file, nil
}

func convertBucketToArr(b *types.Bucket) []string {
	name := b.GetName()
	creationTime := b.GetCreationTime()
	lastModifiedTime := b.GetLastModifiedTime()
	status := b.GetStatus()

	// formatting Time to string
	creationTimeStr := creationTime.Format("2006-01-02T15:04:05-07:00")
	lastModifiedTimeStr := lastModifiedTime.Format("2006-01-02T15:04:05-07:00")

	var data []string

	data = append(data, name, creationTimeStr, lastModifiedTimeStr, status)

	return data
}
