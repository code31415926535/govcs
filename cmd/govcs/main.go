package main

import (
	"fmt"
	"govcs"
	"os"
	"path/filepath"
)

func main() {
	// govcs <command> [args...]
	if len(os.Args) < 2 {
		quit("Command is required.")
	}

	cmd := os.Args[1] // os.Args[0] is govcs

	switch cmd {
	case "init":
		cmdInit(os.Args[2:])
	case "add":
		cmdAdd(os.Args[2:])
	default:
		quit(fmt.Sprintf("Unknown command %s", cmd))
	}
}

func cmdInit(args []string) {
	err := govcs.Init(".")
	if err != nil {
		quit(err.Error())
	}
}

func cmdAdd(args []string) {
	if len(args) != 1 {
		quit("Add takes exactly 1 argument(s)")
	}

	err := govcs.AddFile(ensureAbsolutePath(args[0]))
	if err != nil {
		quit(err.Error())
	}
}

func ensureAbsolutePath(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(getwd(), path)
	}

	return path
}

func getwd() string {
	cdir, err := os.Getwd()
	if err != nil {
		quit(err.Error())
	}

	return cdir
}

func printHelp() {
	fmt.Println("Usage: govcs <command> [args ...]")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  init - initialize a repo.")
	fmt.Println("  add <file> - add file to staging area.")
}

func quit(message string) {
	fmt.Println(message)
	printHelp()
	os.Exit(1)
}
