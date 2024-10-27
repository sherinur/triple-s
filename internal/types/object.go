package types

import "time"

type Object struct {
	ObjectKey        string    `xml:"ObjectKey"`
	Size             string    `xml:"Size"`
	ContentType      string    `xml:"ContentType"`
	LastModifiedTime time.Time `xml:"LastModifiedTime"`
}

func NewObject(objectKey string, size string, contentType string) *Object {
	return &Object{
		ObjectKey:        objectKey,
		Size:             size,
		ContentType:      contentType,
		LastModifiedTime: time.Now(),
	}
}
