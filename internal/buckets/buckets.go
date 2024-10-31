package buckets

import (
	"regexp"
	"strings"
	"time"

	"triple-s/internal/types"
)

// func CreateBucket(name string) error {
// 	bucket := types.NewBucket(name)

// 	csvWriter, file, err := OpenCSV("buckets.csv", true)
// 	if err != nil {
// 		return fmt.Errorf("error of opening or creating a bucket metadata: %w", err)
// 	}
// 	defer func() {
// 		file.Close()
// 	}()

// 	data := ConvertBucketToArr(bucket)
// 	if err := csvWriter.Write(data); err != nil {
// 		return fmt.Errorf("error of writing a bucket metadata: %w", err)
// 	}

// 	csvWriter.Flush()

// 	err = CreateDir("./data/" + name)
// 	if err != nil {
// 		return fmt.Errorf("error of opening or creating a bucket dir: %w", err)
// 	}

// 	return nil
// }

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

func IsUniqueBucketName(name string, records [][]string) bool {
	// TODO: Написать отдельную функцию для проверки уникальности имени бакета
	return true
}

func ValidateBucketName(name string) bool {
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
