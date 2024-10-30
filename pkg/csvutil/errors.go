package csvutil

import "errors"

var (
	ErrNotCSV   = errors.New("the file is not in CSV format")
	ErrNoWriter = errors.New("writer is not initialized")
)
