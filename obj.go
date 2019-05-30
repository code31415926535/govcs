package govcs

import (
	"fmt"
	"govcs/engine"
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
	fmt.Println("Staged changes:")
	for _, f := range s.ChangedFiles {
		fmt.Printf("  %s\n", f)
	}
}
