package types

import "encoding/xml"

type GetBucketResult struct {
	XMLName xml.Name `xml:"GetBucketResult"`
	Bucket  `xml:"Bucket"`
}
