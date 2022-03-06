package querytool

import (
	"errors"
	"strings"
	"testing"

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
	}

	a := assert.New(t)
	for i, test := range tests {
		t.Logf("test #%d", i+1)

		queries, err := loadCSV(strings.NewReader(test.csv))
		a.Nil(queries)
		a.Equal(err, test.err)
	}
}
