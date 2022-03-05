package querytool

import (
	"time"
)

// There's a question whether this time interval should be
// inclusive [start, end] or open-ended [start, end)
// It's also a bit weird to display minute time intervals
// but with start and end times that include seconds.
// I think this query is most true to the requirements.
const cpuStatsQuery = `
	SELECT time_bucket('1 minute', u.ts) as one_min, min(u.usage), max(u.usage)
	FROM cpu_usage u
	WHERE u.host = $1
	AND u.ts BETWEEN $2 AND $3
	GROUP BY one_min
	ORDER BY one_min DESC`

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

	numRows, err := query.executeQuery()
	stats.NumResultRows = numRows
	stats.Duration = time.Now().Sub(start)
	return stats, err
}

func (query *CPUQuery) executeQuery() (int, error) {
	return executeQueryAndDiscardResults(
		cpuStatsQuery, query.Host, query.Start, query.End,
	)
}
