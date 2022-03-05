package main

import (
	"github.com/eloff/queryhw/querytool"
)

func main() {
	options := querytool.ParseCommandOptions()

	querytool.Run(&options)
}
