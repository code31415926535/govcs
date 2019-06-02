package engine

import (
	"bytes"
	"fmt"
	"govcs/metadata"
	"path/filepath"
)

const commitDir = "commits"

// CreateNewCommit with message. The list of diffs and current head are read
// from the index file. The commit is saved but the index file remains unmodified.
//
// The return value is the hash of the commit.
func CreateNewCommit(fs metadata.Metadata, message string) (string, error) {
	i, err := LoadIndex(fs)
	if err != nil {
		return "", err
	}

	c := newCommit(i, message)
	return c.filename(), c.save(fs)
}

// LoadCommits starting from head. maxCount is the maximum number of commits to be
// loaded. If maxCount < 0, all the commits will be loaded.
//
// The commits are returned in reverse order, meaning the most recent
// commit is last and the least recent commit is first.
func LoadCommits(fs metadata.Metadata, maxCount int) ([]*Commit, error) {
	i, err := LoadIndex(fs)
	if err != nil {
		return nil, err
	}

	return loadCommits(fs, i.Head, maxCount)
}

func loadCommits(fs metadata.Metadata, hash string, remaining int) ([]*Commit, error) {
	c, err := loadCommit(fs, hash)
	if err != nil {
		return nil, err
	}

	if remaining == 1 || c.Prev == "" {
		return []*Commit{c}, nil
	}

	partialRes, err := loadCommits(fs, c.Prev, remaining-1)
	if err != nil {
		return nil, err
	}

	return append(partialRes, c), nil
}

func loadCommit(fs metadata.Metadata, hash string) (*Commit, error) {
	d, err := fs.ReadFile(filepath.Join(commitDir, hash))
	if err != nil {
		return nil, err
	}

	parts := bytes.Split(d, []byte("\n"))
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid commit: %s", hash)
	}

	prev := string(parts[0])
	message := string(parts[1])
	diffs := make(map[string]string)

	for _, part := range parts[2:] {
		diff := bytes.Split(part, []byte(" "))
		if len(diff) != 2 {
			return nil, fmt.Errorf("invalid commit: %s", hash)
		}

		diffs[string(diff[0])] = string(diff[1])
	}

	return &Commit{
		Prev:    prev,
		Message: message,
		Diffs:   diffs,
	}, nil
}

func newCommit(i *Index, message string) *Commit {
	return &Commit{
		Prev:    i.Head,
		Message: message,
		Diffs:   i.Diffs,
	}
}

type Commit struct {
	Prev    string
	Message string
	Diffs   map[string]string
}

func (c *Commit) save(fs metadata.Metadata) error {
	return fs.WriteFile(filepath.Join(commitDir, c.filename()), c.data())
}

func (c *Commit) filename() string {
	data := append([]byte(c.Prev), []byte(c.Message)...)
	for k, v := range c.Diffs {
		data = append(data, []byte(k)...)
		data = append(data, []byte(v)...)
	}
	return dataHash(data)
}

func (c *Commit) data() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s\n%s", c.Prev, c.Message)
	for k, v := range c.Diffs {
		fmt.Fprintf(&b, "\n%s %s", k, v)
	}
	return b.Bytes()
}
