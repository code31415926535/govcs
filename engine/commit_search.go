package engine

import "github.com/code31415926535/govcs/metadata"

// FindLatestDiffOfFile in repository. The search begins at head and moves
// down the commit chain until a commit is found that modifies the file.
//
// The return value is the hash of the diff. If no value if found, the return
// value is "" (empty string).
func FindLatestDiffOfFile(fs metadata.Metadata, file string) (string, error) {
	i, err := LoadIndex(fs)
	if err != nil {
		return "", err
	}

	commit, err := findCommitMatching(fs, i.Head, func(c *Commit) bool {
		_, ok := c.Diffs[file]
		return ok
	})

	if commit == nil {
		return "", nil
	}

	return commit.Diffs[file], nil
}

func findCommitMatching(fs metadata.Metadata, hash string, f func(c *Commit) bool) (*Commit, error) {
	// Empty commit, hit rock bottom
	if hash == "" {
		return nil, nil
	}

	c, err := loadCommit(fs, hash)
	if err != nil {
		return nil, err
	}

	if f(c) {
		return c, nil
	}

	return findCommitMatching(fs, c.Prev, f)
}
