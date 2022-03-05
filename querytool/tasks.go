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
	// TODO
	stats := QueryStats{}
	time.Sleep(1 * time.Second)
	return stats, nil
}
