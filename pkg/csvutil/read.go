package csvutil

import (
	"errors"
	"fmt"
)

// ReadAllRecords() reads all records from the CSV-file.
// Returs two-dim slice string and error if it occurs.
func (csvFile *CSVFile) ReadAllRecords() ([][]string, error) {
	if csvFile.Reader == nil {
		return nil, errors.New("reader is not initialized")
	}
	records, err := csvFile.Reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can not read records: %w", err)
	}
	return records, nil
}

// FindInCSV() searches for a record with the specified value in the given slice.
// Returns index of record and true if found, and -1, false if did not.
func FindInSlice(value string, records [][]string) (int, bool) {
	for index, record := range records {
		for _, field := range record {
			if field == value {
				return index, true
			}
		}
	}

	return -1, false
}
