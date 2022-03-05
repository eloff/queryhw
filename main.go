package main

import (
	"runtime/debug"
	"time"

	"github.com/eloff/queryhw/querytool"
)

func main() {
	// Disable the GC. We shouldn't need it for the purposes of this tool.
	// It may make the benchmark results less reliable.
	// A number of Go CLI tools do this.
	debug.SetGCPercent(-1)

	options := querytool.ParseCommandOptions()

	start := time.Now()
	stats := querytool.Run(&options)

	querytool.PrintSummaryStats(&options, time.Now().Sub(start), stats)
}
