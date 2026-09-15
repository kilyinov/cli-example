package main

import (
	"os"

	"github.com/kilyinov/cli-example/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
