package querytool

import (
	"flag"
	"runtime"
)

type Options struct {
	DBConnectionString string
	InputFilePath      string
	NumWorkers         int
	Verbose            bool
}

// This is not a good idea in a real app
// The credentials will be stored in the binary
// where anyone can read them.
const dbConnectionStr = "postgres://postgres:password@db/homework?sslmode=disable"

// ParseCommandOptions parses the CLI options and returns them as an Options struct
func ParseCommandOptions() Options {
	var options Options

	// Define the command line flags that we accept, and their default values
	numWorkers := flag.Int("n", runtime.GOMAXPROCS(0), "the number of concurrent workers to run")
	queriesFile := flag.String("f", "-", "the path to a CSV file containing the queries to run")
	verbose := flag.Bool("v", false, "print more verbose output as the program runs")
	dbConnString := flag.String("d", dbConnectionStr, "database connection string for timescaledb, see docs for lib/pq")

	flag.Parse()

	// Copy the values into the Options struct and do any validation here
	options.NumWorkers = *numWorkers
	options.InputFilePath = *queriesFile
	options.Verbose = *verbose
	options.DBConnectionString = *dbConnString

	return options
}
