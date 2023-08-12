package main

import (
	"os"

	"github.com/kodmain/kitsune/src/cmd/kitsune"
)

func main() {
	if err := kitsune.Helper.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
