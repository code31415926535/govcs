package engine

import (
	"encoding/json"
	"fmt"
	"govcs/metadata"
)

const indexFilePath = "index"

func CreateNewIndex(fs metadata.Metadata) error {
	return newIndex().save(fs)
}

func AddDiffToIndex(fs metadata.Metadata, path, hash string) error {
	i, err := LoadIndex(fs)
	if err != nil {
		return err
	}

	i.Diffs[path] = hash

	return i.save(fs)
}

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

func ChangeHead(fs metadata.Metadata, hash string) error {
	return changeHead(fs, hash, false)
}

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
type Index struct {
	// Head - is the name of the current commit. It can be empty, signaling
	//	that the head refers to no commits.
	Head string `json:"head"`

	// Diffs - is a map containg <filename>:<diff_hash> entries. These are
	//	the entries staged for commit.
	Diffs map[string]string `json:"diffs"`
}

func (i *Index) save(fs metadata.Metadata) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return fs.WriteFile(indexFilePath, d)
}
