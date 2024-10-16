package types

import "time"

type Bucket struct {
	Name string
	CreationTime time.Time
	LastModifiedTime time.Time
	Status string
}