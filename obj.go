package govcs

import (
	"fmt"
	"github.com/code31415926535/govcs/engine"
)

func newStatusFromIndex(i *engine.Index) Status {
	var cf []string
	for k := range i.Diffs {
		cf = append(cf, k)
	}

	return Status{
		Ref:          i.Head,
		ChangedFiles: cf,
	}
}

type Status struct {
	Ref          string
	ChangedFiles []string
}

// TODO: Detach print from obj
func (s Status) Print() {
	fmt.Printf("Ref: %s\n", s.Ref)
	fmt.Println()
	if len(s.ChangedFiles) > 0 {
		fmt.Println("Staged changes:")
		for _, f := range s.ChangedFiles {
			fmt.Printf("  %s\n", f)
		}
	} else {
		fmt.Println("No changes staged.")
	}
}

func newCommit(id string, c *engine.Commit) Commit {
	return Commit{
		ID:      id,
		Prev:    c.Prev,
		Message: c.Message,
	}
}

type Commit struct {
	ID      string
	Prev    string
	Message string
}

// TODO: Detach print from obj
func (c Commit) Print() {
	fmt.Printf("%s -> %s  %s\n", c.ID, c.Prev, c.Message)
}

func newCommits(head string, cs []*engine.Commit) Commits {
	commits := Commits{}

	if len(cs) == 0 {
		return commits
	}

	for id, c := range cs[:len(cs)-1] {
		commits = append(commits, newCommit(cs[id+1].Prev, c))
	}
	commits = append(commits, newCommit(head, cs[len(cs)-1]))

	return commits
}

type Commits []Commit

// TODO: Detach print from obj
func (cs Commits) Print() {
	if len(cs) == 0 {
		fmt.Println("No commits.")
	} else {
		last := len(cs) - 1
		for i := range cs {
			cs[last-i].Print()
		}
	}
}
