package querytool

import (
	"fmt"
	"io"
	"os"
)

// LoadTasks loads t
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
		return nil, fmt.Errorf("LoadTasks error in CSV input: %w", err)
	}

	// TODO group queries by host
	_ = queries
	var tasks []QueryTask

	// TODO sort tasks by number of queries, descending

	return NewTaskQueue(tasks), nil
}

func loadCSV(reader io.Reader) ([]CPUQuery, error) {
	// TODO
	return nil, nil
}
