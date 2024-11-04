package objects

import (
	"regexp"
	"strings"
	"time"

	"triple-s/internal/types"
)

func ValidateObjectKey(objectKey string) bool {
	if len(objectKey) == 0 || len(objectKey) > 1024 {
		return false
	}

	validKeyPattern := regexp.MustCompile(`^[\w.\-\/]+$`)
	if !validKeyPattern.MatchString(objectKey) {
		return false
	}

	if strings.HasPrefix(objectKey, "/") || strings.HasSuffix(objectKey, "/") {
		return false
	}

	return true
}

func ConvertObjectToArr(o *types.Object) []string {
	var data []string

	lastModifiedTimeStr := o.LastModifiedTime.Format("2006-01-02T15:04:05-07:00")
	data = append(data, o.ObjectKey, o.Size, o.ContentType, lastModifiedTimeStr)

	return data
}

func ConvertArrToObject(data []string) (types.Object, error) {
	layout := "2006-01-02T15:04:05-07:00"

	time, err := time.Parse(layout, data[3])
	if err != nil {
		return types.Object{}, err
	}
	return types.Object{
		ObjectKey:        data[0],
		Size:             data[1],
		ContentType:      data[2],
		LastModifiedTime: time,
	}, nil
}
