package csvutil

import (
	"encoding/csv"
	"os"
)

// CSVFile struct to work with files in CSV format
// File is a pointer to the opened file (os.File)
// Writer is a pointer to the writer
type CSVFile struct {
	File   *os.File
	Writer *csv.Writer
	Reader *csv.Reader
}
