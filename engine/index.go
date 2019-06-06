package engine

import (
	"encoding/json"
	"fmt"

	"github.com/code31415926535/govcs/metadata"
)

const indexFilePath = "index"

// CreateNewIndex creates a new empty index file.
func CreateNewIndex(fs metadata.Metadata) error {
	return newIndex().save(fs)
}

// AddDiffTOIndex adds diff of path with hash to index.
func AddDiffToIndex(fs metadata.Metadata, path, hash string) error {
	i, err := LoadIndex(fs)
	if err != nil {
		return err
	}

	i.Diffs[path] = hash

	return i.save(fs)
}

// RemoveDiffFromIndex removes staged diff from index and deletes
// the diff file. Since the diff is stored in index, it means
// it hasn't been commited yet so it is safe to delete the file.
func RemoveDiffFromIndex(fs metadata.Metadata, path string) error {
	i, err := LoadIndex(fs)
	if err != nil {
		return err
	}

	diff := i.Diffs[path]
	// REVIEW - Where should this end up? Here or in vcs
	err = RemoveDiff(fs, diff)
	if err != nil {
		return err
	}
	delete(i.Diffs, path)

	return i.save(fs)
}

// ChangeHead changes the current head to hash. If the
// working directory is not clean, it will return an error.
func ChangeHead(fs metadata.Metadata, hash string) error {
	return changeHead(fs, hash, false)
}

// ChangeHeadForce changes the current head to hash even
// if there are staged changes. All changes will be
// cleared as a result of the change.
func ChangeHeadForce(fs metadata.Metadata, hash string) error {
	return changeHead(fs, hash, true)
}

func changeHead(fs metadata.Metadata, hash string, force bool) error {
	i, err := LoadIndex(fs)
	if err != nil {
		return err
	}

	if !force && len(i.Diffs) != 0 {
		return fmt.Errorf("could not change head due to staged changes")
	}

	i.Diffs = make(map[string]string)
	i.Head = hash

	return i.save(fs)
}

// LoadIndex loads index file.
func LoadIndex(fs metadata.Metadata) (*Index, error) {
	indexData, err := fs.ReadFile(indexFilePath)
	if err != nil {
		return nil, err
	}

	var i Index
	err = json.Unmarshal(indexData, &i)
	return &i, err
}

func newIndex() *Index {
	return &Index{
		Head:  "",
		Diffs: make(map[string]string),
	}
}

// Index contains repository metadata and it is stored in a json format.
// The index file keeps track of the current head and staged changes.
type Index struct {
	// Head is the hash of the current commit. It can be empty, signaling
	// that the head refers to no commit or is the root of the current tree.
	Head string `json:"head"`

	// Diffs is a map containg <filename>:<diff_hash> entries. These are
	// the entries staged for commit.
	Diffs map[string]string `json:"diffs"`
}

func (i *Index) save(fs metadata.Metadata) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return fs.WriteFile(indexFilePath, d)
}
