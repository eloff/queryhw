package querytool

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadCSVError(t *testing.T) {
	tests := []struct {
		csv string
		err error
	}{
		// Malformed header
		{
			csv: "host,start,end\n",
			err: errors.New("start time must be formatted like 2006-01-02 15:04:05, not start"),
		},
		// Too many values
		{
			csv: "12,4,5,6",
			err: errors.New("expected CSV row to contain 3 values: got 4"),
		},
		// Too few values
		{
			csv: "12,4",
			err: errors.New("expected CSV row to contain 3 values: got 2"),
		},
		// Malformed start date (missing time)
		{
			csv: "1,2006-01-02,2006-01-02 15:04:05",
			err: errors.New("start time must be formatted like 2006-01-02 15:04:05, not 2006-01-02"),
		},
		// Malformed start date (invalid)
		{
			csv: "bar,2006-13-02 15:04:05,2006-01-02 15:04:05",
			err: errors.New("start time must be formatted like 2006-01-02 15:04:05, not 2006-13-02 15:04:05"),
		},
		// Malformed end date (time zone)
		{
			csv: "hostname,start,end\nfoo,2006-01-02 15:04:05,2006-01-02 15:04:05Z07:00",
			err: errors.New("end time must be formatted like 2006-01-02 15:04:05, not 2006-01-02 15:04:05Z07:00"),
		},
		// Not a CSV
		{
			csv: "{foo\tbar\tbaz}",
			err: errors.New("expected CSV row to contain 3 values: got 1"),
		},
	}

	a := assert.New(t)
	for i, test := range tests {
		t.Logf("test #%d", i+1)

		queries, err := loadCSV(strings.NewReader(test.csv))
		a.Nil(queries)
		a.Equal(err, test.err)
	}
}

func TestLoadCSV(t *testing.T) {
	csv := `hostname,start,end
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28
host_000002,2017-01-02 15:16:29,2017-01-02 16:16:29
host_000003,2017-01-01 08:52:14,2017-01-01 09:52:14
`
	expected := []CPUQuery{
		{
			Host:  "host_000008",
			Start: time.Date(2017, 1, 1, 8, 59, 22, 0, time.UTC),
			End:   time.Date(2017, 1, 1, 9, 59, 22, 0, time.UTC),
		},
		{
			Host:  "host_000001",
			Start: time.Date(2017, 1, 2, 13, 2, 2, 0, time.UTC),
			End:   time.Date(2017, 1, 2, 14, 2, 2, 0, time.UTC),
		},
		{
			Host:  "host_000008",
			Start: time.Date(2017, 1, 2, 18, 50, 28, 0, time.UTC),
			End:   time.Date(2017, 1, 2, 19, 50, 28, 0, time.UTC),
		},
		{
			Host:  "host_000002",
			Start: time.Date(2017, 1, 2, 15, 16, 29, 0, time.UTC),
			End:   time.Date(2017, 1, 2, 16, 16, 29, 0, time.UTC),
		},
		{
			Host:  "host_000003",
			Start: time.Date(2017, 1, 1, 8, 52, 14, 0, time.UTC),
			End:   time.Date(2017, 1, 1, 9, 52, 14, 0, time.UTC),
		},
	}

	queries, err := loadCSV(strings.NewReader(csv))

	assert.Nil(t, err)
	assert.Equal(t, queries, expected)
}
