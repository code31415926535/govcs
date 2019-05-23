package main

import (
	"govcs"
	"os"
)

func main() {
	// govcs <command> [args...]
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	cmd := os.Args[1] // os.Args[0] is govcs

	switch cmd {
	case "init":
		cmdInit()
	}
}

func cmdInit() {
	err := govcs.NewRepository(".")
	if err != nil {
		os.Exit(1)
	}
}
