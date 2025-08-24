package main

import (
	"os"

	"github.com/alowayed/go-univers/cmd/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}

// testing
