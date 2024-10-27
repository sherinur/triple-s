package types

import "time"

type Bucket struct {
	Name             string    `xml:"Name"`
	CreationTime     time.Time `xml:"CreationTime"`
	LastModifiedTime time.Time `xml:"LastModifiedTime"`
	Status           string    `xml:"Status"`
}

func NewBucket(name string) *Bucket {
	return &Bucket{
		Name:             name,
		CreationTime:     time.Now(),
		LastModifiedTime: time.Now(),
		Status:           "active",
	}
}
