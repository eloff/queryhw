package querytool

import (
	"fmt"
	"sort"
	"time"
)

// QueryStats contains benchmark stats from running a query
type QueryStats struct {
	WorkerId      int
	NumResultRows int
	Duration      time.Duration
	Host          string
}

// IsZero returns true if this QueryStats struct is zero initialized
func (stats *QueryStats) IsZero() bool {
	return stats.WorkerId == 0 && stats.Duration == 0 && stats.Host == ""
}

// ByNumberOfQueries implements sort.Interface for []QueryTask based on the number of queries (descending)
type ByDurationDescending []QueryStats

func (a ByDurationDescending) Len() int           { return len(a) }
func (a ByDurationDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDurationDescending) Less(i, j int) bool { return a[i].Duration > a[j].Duration }

// PrintSummaryStats prints the summary statistics for all the queries run
func PrintSummaryStats(options *Options, totalDuration time.Duration, allStats []QueryStats) {
	stats := calculateSummaryStats(allStats)

	fmt.Printf("Executed %d queries in %.2f seconds\n", len(allStats), float64(totalDuration)/float64(time.Second))
	fmt.Printf("Total execution time for all queries was %.2f seconds, using %d worker threads. Parallel speedup of %.1fx\n",
		float64(stats.Total)/float64(time.Second), options.NumWorkers, float64(stats.Total)/float64(totalDuration))
	fmt.Printf(`min query duration=%.2fms
max query duration=%.2fms
average=%.2fms
median=%.2fms
95th percentile=%.2fms
standard deviation=%.2fms
`,
		float64(stats.Min)/float64(time.Millisecond),
		float64(stats.Max)/float64(time.Millisecond),
		float64(stats.Average)/float64(time.Millisecond),
		float64(stats.Median)/float64(time.Millisecond),
		.0, // TODO
		.0, // TODO
	)
}

type SummaryStats struct {
	// We can use Duration (int64) for median and average
	// because we don't need fractional nanoseconds.
	// Computer clocks are just not that accurate,
	// and neither is our benchmark code. We're actually going
	// to drop the nanoseconds and microseconds when we display it anyway.
	Min, Max, Total, Median, Average time.Duration
}

// calculateSummaryStats computes the summary statistics for all the queries.
// We use a separate method because we want to write unit tests for it.
func calculateSummaryStats(allStats []QueryStats) SummaryStats {
	if len(allStats) == 0 {
		// This is programmer error, not a runtime error, so we panic
		panic("allStats cannot be empty")
	}

	// This is not ideal sorting an array of structs since it means a lot of copying
	// We need this to compute percentile values like the median or 95th percentile.
	sort.Sort(ByDurationDescending(allStats))

	var min, max, total, average, median time.Duration
	min = allStats[len(allStats)-1].Duration
	max = allStats[0].Duration
	for _, stats := range allStats {
		total += stats.Duration
	}

	// Compute the median and the average
	mid := len(allStats) / 2
	median = allStats[mid].Duration
	if mid*2 == len(allStats) {
		// We have an even number of queries, to compute the median
		// take the two middle values and divide by 2.
		median = (median + allStats[mid-1].Duration) / 2
	}
	average = total / time.Duration(len(allStats))

	return SummaryStats{
		Min:     min,
		Max:     max,
		Total:   total,
		Median:  median,
		Average: average,
	}
}
