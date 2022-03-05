package querytool

import (
	"log"
)

func Run(options *Options) {
	tasks, err := LoadTasks(options.InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	results := make(chan QueryStats, tasks.Len())

	// Launch the workers
	for i := 0; i < options.NumWorkers; i++ {
		go runWorker(i, tasks, results)
	}
}

// runWorker runs a worker goroutine that will process tasks
// from the TaskQueue one at a time, sending the results to
// the main goroutine via the results channel.
func runWorker(id int, tasks *TaskQueue, results chan QueryStats) {
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
			results <- stats
		}
	}
}
