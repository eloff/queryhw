package querytool

import (
	"fmt"
	"math"
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

// ByDuration implements sort.Interface for []QueryStats based on the Duration (ascending)
type ByDuration []QueryStats

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i].Duration < a[j].Duration }

// PrintSummaryStats prints the summary statistics for all the queries run
func PrintSummaryStats(options *Options, totalDuration time.Duration, allStats []QueryStats) {
	stats := calculateSummaryStats(allStats)

	// I asked about the purpose of this program and who the users might be.
	// I was told not to worry about it. So I'm very much guessing here what
	// summary statistics might be interesting to the user.
	//
	// Since it's a benchmark program, I output some stats about
	// how long it took, how many queries were executed, what
	// parallel speedup was acheived by using the specified number of workers.
	//
	// Since the requirements specify a min, max, median, and average
	// value, I add the 95th percentile and standard deviation as those
	// may also be interesting to the user.

	// Print how many queries we executed and the "walltime" elapsed
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
		float64(stats._95Percentile)/float64(time.Millisecond), // TODO
		stats.StdDev, // TODO
	)
}

type SummaryStats struct {
	Min, Max, Total, Median, Average, _95Percentile time.Duration
	StdDev                                          float64
}

// calculateSummaryStats computes the summary statistics for all the queries.
// We use a separate method because we want to write unit tests for it.
func calculateSummaryStats(allStats []QueryStats) SummaryStats {
	if len(allStats) == 0 {
		// This is programmer error, not a runtime error, so we panic
		panic("allStats cannot be empty")
	}

	// This is not ideal sorting an array of structs since it means a lot of copying.
	// Sorting an array of pointers to structs also has its issues (cache misses, allocations.)
	// We need this to compute percentile values like the median or 95th percentile.
	sort.Sort(ByDuration(allStats))

	// We can use Duration (int64) for median and average
	// because we don't need fractional nanoseconds.
	// Computer clocks are just not that accurate,
	// and neither is our benchmark code. We're actually going
	// to drop the nanoseconds and microseconds when we display it anyway.
	var min, max, total, average, median, _95percentile time.Duration
	min = allStats[0].Duration
	max = allStats[len(allStats)-1].Duration

	// We'll use Welford's algorithm for estimation.
	// This algorithm is more numerically stable than most other algorithms
	// for calculating the variance.
	n := .0
	mu := .0
	sq := .0

	for _, stats := range allStats {
		total += stats.Duration

		n++
		x := float64(stats.Duration) / float64(time.Millisecond)
		nextMu := mu + (x-mu)/n
		sq += (x - mu) * (x - nextMu)
		mu = nextMu
	}
	stddev := math.Sqrt(sq / n)

	// Compute the median and the average
	mid := len(allStats) / 2

	// There are multiple ways to compute the median (50th percentile).
	// Since we also compute the 95th percentile we'll use
	// the same algorithm for both. It shouldn't matter much
	// and it's easy enough to change if we must.
	median = computePercentile(allStats, 0.5)
	_95percentile = computePercentile(allStats, 0.95)

	median = allStats[mid].Duration
	if mid*2 == len(allStats) {
		// We have an even number of queries, to compute the median
		// take the two middle values and divide by 2.
		median = (median + allStats[mid-1].Duration) / 2
	}
	average = total / time.Duration(len(allStats))

	return SummaryStats{
		Min:           min,
		Max:           max,
		Total:         total,
		Median:        median,
		Average:       average,
		_95Percentile: _95percentile,
		StdDev:        stddev,
	}
}

func computePercentile(allStats []QueryStats, percentile float64) time.Duration {
	index := float64(len(allStats)) * percentile

	// If the index is a whole number, then the percentile value
	// is the average of the value at the index and the value that follows
	// Otherwise we round it up to the nearest whole number and
	// the value is the value at that index.

	roundedUp := math.Ceil(index)
	value := allStats[int(roundedUp)-1].Duration
	// I'm not 100% sure this shouldn't be abs(roundedUp-index) < epsilon
	if roundedUp == index {
		// Index was a whole number so take the average
		value = (value + allStats[int(roundedUp)].Duration) / 2
	}

	return value
}
