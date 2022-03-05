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

func (query *CPUQuery) Run() (QueryStats, error) {
	// TODO
	stats := QueryStats{}
	time.Sleep(1 * time.Second)
	return stats, nil
}
