package engine

import (
	"bytes"
	"fmt"
	"govcs/metadata"
	"path/filepath"
)

func CreateNewCommit(fs metadata.Metadata, message string) (string, error) {
	i, err := LoadIndex(fs)
	if err != nil {
		return "", err
	}

	c := newCommit(i, message)
	return c.filename(), c.save(fs)
}

func newCommit(i *Index, message string) *commit {
	return &commit{
		Prev:    i.Head,
		Message: message,
		Diffs:   i.Diffs,
	}
}

type commit struct {
	Prev    string
	Message string
	Diffs   map[string]string
}

func (c *commit) save(fs metadata.Metadata) error {
	return fs.WriteFile(filepath.Join("commits", c.filename()), c.data())
}

func (c *commit) filename() string {
	data := append([]byte(c.Prev), []byte(c.Message)...)
	for k, v := range c.Diffs {
		data = append(data, []byte(k)...)
		data = append(data, []byte(v)...)
	}
	return dataHash(data)
}

func (c *commit) data() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s\n%s", c.Prev, c.Message)
	for k, v := range c.Diffs {
		fmt.Fprintf(&b, "\n%s %s", k, v)
	}
	return b.Bytes()
}
