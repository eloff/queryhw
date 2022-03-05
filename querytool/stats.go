package querytool

import (
	"fmt"
	"math"
	"time"
)

// QueryStats contains benchmark stats from running a query
type QueryStats struct {
	WorkerId int
	Duration time.Duration
	Host     string
}

// IsZero returns true if this QueryStats struct is zero initialized
func (stats *QueryStats) IsZero() bool {
	return stats.WorkerId == 0 && stats.Duration == 0 && stats.Host == ""
}

// PrintSummaryStats prints the summary statistics for all the queries run
func PrintSummaryStats(options *Options, totalDuration time.Duration, allStats []QueryStats) {
	stats := calculateSummaryStats(allStats)

	fmt.Printf("Executed %d queries in %.2f seconds\n", len(allStats), float64(totalDuration)/float64(time.Second))
	fmt.Printf("Total execution time for all queries was %.2f seconds, using %d worker threads. Parallel speedup of %.1fx\n",
		float64(stats.Total)/float64(time.Second), options.NumWorkers, float64(stats.Total)/float64(totalDuration))
	fmt.Printf("min=%.2fms, max=%.2fms, median=%.2fms, average=%.2fms\n",
		float64(stats.Min)/float64(time.Millisecond),
		float64(stats.Max)/float64(time.Millisecond),
		float64(stats.Median)/float64(time.Millisecond),
		float64(stats.Average)/float64(time.Millisecond))
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
	var min, max, total, average, median time.Duration
	min = time.Duration(math.MaxInt64)
	for _, stats := range allStats {
		if min > stats.Duration {
			min = stats.Duration
		}
		if max < stats.Duration {
			max = stats.Duration
		}
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
