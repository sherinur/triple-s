package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"triple-s/internal/types"
)

// func CreateObjectMeta(name string) error {
// 	object := types.NewObject(name)
// }

func CreateBucket(name string) error {
	bucket := types.NewBucket(name)

	csvWriter, file, err := openCSV("buckets.csv")
	if err != nil {
		return fmt.Errorf("error of opening or creating a bucket metadata: %w", err)
	}
	defer func() {
		file.Close()
	}()

	data := ConvertBucketToArr(bucket)
	if err := csvWriter.Write(data); err != nil {
		return fmt.Errorf("error of writing a bucket metadata: %w", err)
	}

	csvWriter.Flush()

	err = CreateDir("./data/" + name)
	if err != nil {
		return fmt.Errorf("error of opening or creating a bucket dir: %w", err)
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

func ParseCSV(filePath string) ([][]string, error) {
	dirPath := filepath.Dir(filePath)

	if err := CreateDir(dirPath); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть или создать файл: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV: %w", err)
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

func ConvertBucketToArr(b *types.Bucket) []string {
	// formatting Time to string
	creationTimeStr := b.CreationTime.Format("2006-01-02T15:04:05-07:00")
	lastModifiedTimeStr := b.LastModifiedTime.Format("2006-01-02T15:04:05-07:00")

	var data []string

	data = append(data, b.Name, creationTimeStr, lastModifiedTimeStr, b.Status)

	return data
}

func ConvertArrToBucket(data []string) (types.Bucket, error) {
	layout := "2006-01-02T15:04:05-07:00"

	t1, err := time.Parse(layout, data[1])
	if err != nil {
		return types.Bucket{}, err
	}

	t2, err := time.Parse(layout, data[2])
	if err != nil {
		return types.Bucket{}, err
	}

	return types.Bucket{
		Name:             data[0],
		CreationTime:     t1,
		LastModifiedTime: t2,
		Status:           data[3],
	}, nil
}
