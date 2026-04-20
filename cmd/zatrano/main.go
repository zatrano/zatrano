package main

import (
	"os"

	"github.com/zatrano/zatrano/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
