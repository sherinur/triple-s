package types

import "time"

type Bucket struct {
	name             string
	creationTime     time.Time
	lastModifiedTime time.Time
	status           string
}

func NewBucket(name string) *Bucket {
	return &Bucket{
		name:             name,
		creationTime:     time.Now(),
		lastModifiedTime: time.Now(),
		status:           "active",
	}
}

func (b *Bucket) GetName() string {
	return b.name
}

func (b *Bucket) GetCreationTime() time.Time {
	return b.creationTime
}

func (b *Bucket) GetLastModifiedTime() time.Time {
	return b.lastModifiedTime
}

func (b *Bucket) GetStatus() string {
	return b.status
}

func (b *Bucket) SetName(name string) {
	b.name = name
	b.lastModifiedTime = time.Now()
}

func (b *Bucket) SetStatus(status string) {
	b.status = status
	b.lastModifiedTime = time.Now()
}
