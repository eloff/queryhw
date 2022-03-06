package querytool

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

	tasks := make([]QueryTask, 0, len(queries))
	if len(queries) == 0 {
		log.Fatal("No input queries given")
	}

	// Group queries by hostname, we can do this by grouping with a hash map
	// or by sorting, since the order between queries doesn't matter.
	// A hash map is more efficient: O(N) vs O(N*Log(N)), but we'll need
	// multiple allocations with the hash map (for map and subarrays)
	// while the sort is in-place.
	//
	// There's no easy way to know which is best, so I'll just take a guess.
	groupedQueries := make(map[string]QueryTask, len(queries))
	for _, query := range queries {
		task := groupedQueries[query.Host]
		task.Queries = append(task.Queries, query)
		groupedQueries[query.Host] = task
	}

	// Collect all the tasks into a slice
	for _, task := range groupedQueries {
		tasks = append(tasks, task)
	}

	// Sort tasks by number of queries, descending
	// This is not necessary, but it seems to result in slightly
	// better CPU utilization under some workloads.
	// Otherwise a worker can end up processing a large task with many queries
	// at the end while the other workers are idle. It's better to do the big
	// tasks first.
	sort.Sort(ByNumberOfQueries(tasks))

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
			return nil, fmt.Errorf("expected CSV row to contain 3 values: got %d", len(record))
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

		// There's the question of timezones here.
		// The database uses timestamptz and the input
		// data in cpu_usage.csv doesn't specify the timezone.
		// PostgreSQL handles this by storing everything as UTC.
		// The input times in query_params.csv that we're loading
		// here also don't have timezones. That means we need to
		// treat these as UTC values as well so both the
		// input and query times are treated the same way.
		// Go also assumes UTC, so we don't need any code here.

		query := CPUQuery{
			Host:  record[0],
			Start: start,
			End:   end,
		}
		queries = append(queries, query)
	}

	return queries, nil
}
