package engine

import (
	"encoding/json"
	"govcs/metadata"
)

const indexFilePath = "index"

func CreateNewIndex(fs metadata.Metadata) error {
	return newIndex().save(fs)
}

func AddDiffToIndex(fs metadata.Metadata, path, hash string) error {
	i, err := loadIndex(fs)
	if err != nil {
		return err
	}

	i.Diffs[path] = hash

	return i.save(fs)
}

func newIndex() *index {
	return &index{
		Head:  "",
		Diffs: make(map[string]string),
	}
}

func loadIndex(fs metadata.Metadata) (*index, error) {
	indexData, err := fs.ReadFile(indexFilePath)
	if err != nil {
		return nil, err
	}

	var i index
	err = json.Unmarshal(indexData, &i)
	return &i, err
}

// Index contains repository metadata and it is stored in a json format.
type index struct {
	// Head - is the name of the current commit. It can be empty, signaling
	//	that the head refers to no commits.
	Head string `json:"head"`

	// Diffs - is a map containg <filename>:<diff_hash> entries. These are
	//	the entries staged for commit.
	Diffs map[string]string `json:"diffs"`
}

func (i *index) save(fs metadata.Metadata) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return fs.WriteFile(indexFilePath, d)
}
