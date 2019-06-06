package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/code31415926535/govcs"
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
	case "remove":
		cmdRemove(os.Args[2:])
	case "stat":
		cmdStat(os.Args[2:])
	case "commit":
		cmdCommit(os.Args[2:])
	case "list-commits":
		cmdListCommits(os.Args[2:])
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

func cmdRemove(args []string) {
	if len(args) != 1 {
		quit("Remove takes exactly 1 argument(s)")
	}

	err := govcs.RemoveFile(ensureAbsolutePath(args[0]))
	if err != nil {
		quit(err.Error())
	}
}

func cmdStat(args []string) {
	stat, err := govcs.Stat(".")
	if err != nil {
		quit(err.Error())
	}

	stat.Print()
}

func cmdCommit(args []string) {
	if len(args) != 1 {
		quit("Commit takes exactly 1 argument(s)")
	}

	err := govcs.CommitChanges(".", args[0])
	if err != nil {
		quit(err.Error())
	}
}

func cmdListCommits(args []string) {
	commits, err := govcs.ListCommits(".")
	if err != nil {
		quit(err.Error())
	}

	commits.Print()
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
	fmt.Println("  remove <file> - remove file from staging area.")
	fmt.Println("  stat - print current head and staging area.")
	fmt.Println("  commit <message> - commit changes.")
	fmt.Println("  list-commits - list commits.")
}

func quit(message string) {
	fmt.Println(message)
	printHelp()
	os.Exit(1)
}
