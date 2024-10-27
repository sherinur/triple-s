package types

import (
	"encoding/xml"
)

type ListAllBucketsResult struct {
	XMLName xml.Name `xml:"ListAllBucketsResult"`
	Buckets struct {
		Bucket []Bucket `xml:"Bucket"`
	} `xml:"Buckets"`
}
