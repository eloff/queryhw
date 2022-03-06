package querytool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPercentiles(t *testing.T) {
	data := []int{43, 54, 56, 61, 62, 66, 68, 69, 69, 70, 71, 72, 77, 78, 79, 85, 87, 88, 89, 93, 95, 96, 98, 99, 99}
	stats := make([]QueryStats, len(data))
	for i, num := range data {
		stats[i] = QueryStats{Duration: time.Duration(num)}
	}

	tests := []struct {
		percentile float64
		expected   int
	}{
		{
			percentile: 0.9,
			expected:   98,
		},
		{
			percentile: 0.2,
			expected:   64,
		},
		{
			percentile: 0.5,
			expected:   77,
		},
	}

	a := assert.New(t)
	for i, test := range tests {
		t.Logf("test #%d", i+1)
		value := computePercentile(stats, test.percentile)
		a.Equal(int(value), test.expected)
	}
}

func TestSummaryStats(t *testing.T) {
	data := []int{56354, 34453, 896789, 54362, 425467, 87665, 123413, 356346, 986878, 131374, 97987, 85644}
	stats := make([]QueryStats, len(data))
	for i, num := range data {
		stats[i] = QueryStats{Duration: time.Duration(num) * time.Millisecond}
	}

	a := assert.New(t)
	summary := calculateSummaryStats(stats)
	a.Equal(int(summary.Min), 34453*int(time.Millisecond))
	a.Equal(int(summary.Max), 986878*int(time.Millisecond))
	a.Equal(int(summary.Total), 3336732*int(time.Millisecond))
	a.Equal(int(summary._95Percentile), 986878*int(time.Millisecond))
	a.Equal(int(summary.Average), 278061*int(time.Millisecond))
	a.Equal(int(summary.Median), 110700*int(time.Millisecond))
	a.Equal(summary.StdDev, 319214.8839603713)
}
