package main

import (
	"os"

	"github.com/mathiasdonoso/hook-replay/internal/cli"
)

func main() {
	if err := cli.NewHR(os.Args[1:]).Execute(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
