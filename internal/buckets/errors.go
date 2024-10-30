package buckets

import "errors"

var (
	ErrInvalidName error = errors.New("bucket name is not valid")
	ErrNotUnique   error = errors.New("bucket name is not unique")
)
