package querytool

import (
	"database/sql"

	_ "github.com/lib/pq" // load the Postgres driver
)

var pool *sql.DB

func InitDB(connectionString string) error {
	var err error
	pool, err = sql.Open("postgres", connectionString)
	return err
}

// executeQueryAndDiscardResults runs the query and returns the number of result rows and any error
// This function fetches and discards the result rows.
func executeQueryAndDiscardResults(query string, args ...interface{}) (int, error) {
	rows, err := pool.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	numRows := 0
	for rows.Next() {
		// Discard the result rows
		numRows++
	}

	return numRows, rows.Err()
}
