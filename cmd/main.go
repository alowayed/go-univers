package main

import (
	"os"
)

func main() {
	os.Exit(run(os.Stdout, os.Args[1:]))
}
