package csvutil

import "fmt"

// RecordsToCSV() opens a file in overwrite mode and writes the records to it.
func (csvFile *CSVFile) RecordsToCSV(records [][]string) error {
	if csvFile.Writer == nil {
		return ErrNoWriter
	}
	if err := csvFile.Writer.WriteAll(records); err != nil {
		return fmt.Errorf("can not write records to CSV file: %w", err)
	}
	return nil
}

// appendToCSV() appends record to the end of the CSV-file.
// Takes path of file and slice of strings as argument.
// If there is an error, returns error.
func (csvFile *CSVFile) AppendToCSV(record []string) error {
	if csvFile.Writer == nil {
		return ErrNoWriter
	}
	if err := csvFile.Writer.Write(record); err != nil {
		return fmt.Errorf("can not append record to CSV file: %w", err)
	}

	csvFile.Writer.Flush()

	if err := csvFile.Writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return nil
}
