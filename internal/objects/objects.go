package objects

import (
	"time"

	"triple-s/internal/types"
)

func CreateObject(bucketName, objectKey, size, contentType string) error {
	// object := types.NewObject(objectKey, size, contentType)

	// // Проверяем и создаем директорию для объекта, если она отсутствует
	// objectDirPath := "./data/" + bucketName
	// if err := os.MkdirAll(objectDirPath, os.ModePerm); err != nil {
	// 	return fmt.Errorf("error creating object directory: %w", err)
	// }

	// csvWriter, file, err := OpenCSV(objectDirPath+"/objects.csv", true)
	// if err != nil {
	// 	return fmt.Errorf("error of opening or creating a bucket metadata: %w", err)
	// }
	// defer func() {
	// 	file.Close()
	// }()

	// data := ConvertObjectToArr(object)
	// if err := csvWriter.Write(data); err != nil {
	// 	return fmt.Errorf("error of writing an object metadata: %w", err)
	// }

	// // Сбрасываем данные для записи в файл
	// csvWriter.Flush()

	return nil
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
