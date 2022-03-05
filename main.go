package main

import (
	"runtime/debug"

	"github.com/eloff/queryhw/querytool"
)

func main() {
	// Disable the GC. We shouldn't need it for the purposes of this tool.
	// It may also make the benchmark results less reliable.
	debug.SetGCPercent(-1)

	options := querytool.ParseCommandOptions()

	querytool.Run(&options)
}
