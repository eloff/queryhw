package querytool

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

// The supported datetime format in the CSV file
const timeFormat = "2006-01-02 15:04:05"

// LoadTasks loads
func LoadTasks(csvFilePath string) (*TaskQueue, error) {
	input := os.Stdin
	if csvFilePath != "" && csvFilePath != "-" {
		// Load the CSV input from the file at path
		var err error
		input, err = os.Open(csvFilePath)
		if err != nil {
			return nil, fmt.Errorf("LoadTasks failed to open %s: %w", csvFilePath, err)
		}
	}

	queries, err := loadCSV(input)
	if err != nil {
		return nil, fmt.Errorf("LoadTasks: %w", err)
	}

	// TODO group queries by host
	_ = queries
	var tasks []QueryTask

	// TODO sort tasks by number of queries, descending

	return NewTaskQueue(tasks), nil
}

// loadCSV parses the CSV file in reader into a slice of CPUQuery structs
func loadCSV(reader io.Reader) ([]CPUQuery, error) {
	var queries []CPUQuery

	csvReader := csv.NewReader(reader)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}
		if len(record) != 3 {
			return nil, fmt.Errorf("expected CSV row to contain 3 values: %d", len(record))
		}
		if len(queries) == 0 && record[0] == "hostname" {
			// This is the header row, skip it
			continue
		}

		start, err := time.Parse(timeFormat, record[1])
		if err != nil {
			return nil, fmt.Errorf("start time must be formatted like %s, not %s", timeFormat, record[1])
		}
		end, err := time.Parse(timeFormat, record[2])
		if err != nil {
			return nil, fmt.Errorf("end time must be formatted like %s, not %s", timeFormat, record[2])
		}

		query := CPUQuery{
			Host:  record[0],
			Start: start,
			End:   end,
		}
		queries = append(queries, query)
	}

	return queries, nil
}
