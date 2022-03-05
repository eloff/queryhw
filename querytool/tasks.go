package querytool

import (
	"time"
)

type QueryTask struct {
	Queries []CPUQuery
}

type CPUQuery struct {
	Host  string
	Start time.Time
	End   time.Time
}

// ByNumberOfQueries implements sort.Interface for []QueryTask based on the number of queries (descending)
type ByNumberOfQueries []QueryTask

func (a ByNumberOfQueries) Len() int           { return len(a) }
func (a ByNumberOfQueries) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumberOfQueries) Less(i, j int) bool { return len(a[i].Queries) > len(a[j].Queries) }

func (query *CPUQuery) Run() (QueryStats, error) {
	// The OS and Go can both interrupt this routine, messing up the timing values
	// I'm not going to do this here, but we can disable preemptive
	// goroutine switching for this goroutine (the GC is disabled anyway.)
	// We can also set the thread priority to a higher level for the OS.

	// Go's time struct contains both "wall time" and a monotonic clock.
	// Subtracting two time values will use the monotonic clock value,
	// which is what we want to get an accurate duration calculation.
	start := time.Now()
	stats := QueryStats{Host: query.Host}
	time.Sleep(100 * time.Millisecond)
	stats.Duration = time.Now().Sub(start)
	return stats, nil
}
