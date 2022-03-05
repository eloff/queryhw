package querytool

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

func Run(options *Options) []QueryStats {
	tasks, err := LoadTasks(options.InputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = InitDB(options.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan QueryStats, tasks.Len())
	// liveWorkers is a shared atomic counter decremented when a worker exits
	// when all workers exit then we'll close the results channel and
	// compute the summary statistics below.
	liveWorkers := int32(options.NumWorkers)
	// Launch the workers
	for i := 0; i < options.NumWorkers; i++ {
		go runWorker(i+1, &liveWorkers, tasks, results)
	}

	allStats := make([]QueryStats, 0, tasks.Len())
	for stats := range results {
		if stats.IsZero() {
			// All workers have exited, there will be no new stats
			break
		}
		allStats = append(allStats, stats)
		if options.Verbose {
			fmt.Printf("query for host %s executed in %.2fms by worker %d\n",
				stats.Host, float64(stats.Duration)/float64(time.Millisecond), stats.WorkerId)
		}
	}

	return allStats
}

// runWorker runs a worker goroutine that will process tasks
// from the TaskQueue one at a time, sending the results to
// the main goroutine via the results channel.
func runWorker(
	id int, liveWorkers *int32,
	tasks *TaskQueue, results chan QueryStats) {
	for {
		task := tasks.Get()
		if task == nil {
			break
		}
		for _, query := range task.Queries {
			stats, err := query.Run()
			if err != nil {
				// In a different context we'd want to report this error
				// back to the main thread and maybe gracefully shutdown
				// the workers before exiting, or something else.
				// Here an error means the database is not available
				// or set up correctly (we validated the tasks
				// when loading them) , or that we have a bug.
				// Either way just reporting the error and exiting is
				// the right thing to do, and also the easy solution.

				// If this leaves zombie connections in the PostgreSQL server
				// we'll need to clean up the connections first.
				// I don't think it will though.
				log.Fatalf("error running query: %v", err)
			}
			stats.WorkerId = id
			results <- stats
		}
	}

	// This worker is finished and will exit now
	if atomic.AddInt32(liveWorkers, -1) == 0 {
		// This is the last worker to exit, close the results channel.
		// This will unblock the main goroutine and signal
		// that it can compute the summary statistics.
		close(results)
	}
}
