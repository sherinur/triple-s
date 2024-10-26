package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

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

	// Записываем данные в CSV
	data := convertBucketToArr(bucket)
	if err := csvWriter.Write(data); err != nil {
		return fmt.Errorf("error of writing a bucket metadata: %w", err)
	}

	// Сбрасываем буфер в файл
	csvWriter.Flush()

	// Создаем директорию для bucket
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

	// Создаем директорию, если её нет
	if err := CreateDir(dirPath); err != nil {
		return nil, err
	}

	// Открываем или создаем CSV файл
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть или создать файл: %w", err)
	}
	defer file.Close()

	// Читаем записи CSV
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

	// Определяем путь к директории data
	dataDirPath := filepath.Join(filepath.Dir(execPath), "data")
	if err := os.MkdirAll(dataDirPath, os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("error of dir creation 'data': %w", err)
	}

	// Путь к файлу buckets.csv
	filePath := filepath.Join(dataDirPath, name)

	// Открываем файл для добавления данных
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("error of opening or creating file - 'buckets.csv': %w", err)
	}

	// Буферизация для повышения производительности записи
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
