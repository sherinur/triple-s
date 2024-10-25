package types

import "time"

type Object struct {
	objectKey        string    `xml:"ObjectKey"`
	size             string    `xml:"Size"`
	contentType      string    `xml:"ContentType"`
	lastModifiedTime time.Time `xml:"LastModifiedTime"`
}

func (o *Object) GetObjectKey() string {
	return o.objectKey
}

func (o *Object) GetSize() string {
	return o.size
}

func (o *Object) GetContentType() string {
	return o.contentType
}

func (o *Object) GetLastModifiedTime() time.Time {
	return o.lastModifiedTime
}

func (o *Object) SetObjectKey(key string) {
	o.objectKey = key
}

func (o *Object) SetSize(size string) {
	o.size = size
}

func (o *Object) SetContentType(contentType string) {
	o.contentType = contentType
}

func (o *Object) SetLastModifiedTime(t time.Time) {
	o.lastModifiedTime = t
}
