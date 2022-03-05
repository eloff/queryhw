package querytool

import (
	"time"
)

/// QueryStats contains benchmark stats from running a query
type QueryStats struct {
	Duration time.Duration
}
