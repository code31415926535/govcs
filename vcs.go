package govcs

import (
	"fmt"
	"govcs/engine"
	"govcs/metadata"
	"io/ioutil"
	"path/filepath"
)

func NewDefaultVcs(path string) (Vcs, error) {
	absPath, err := filepath.Abs(path)
	return Vcs{
		fs:   metadata.NewFileSystemMetadata(absPath),
		root: absPath,
	}, err
}

func LoadDefaultVcs(path string) (Vcs, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Vcs{}, err
	}

	return loadDefaultVcs(absPath)
}

func loadDefaultVcs(path string) (Vcs, error) {
	parent := filepath.Dir(path)
	// Hit root of filesystem
	if parent == path {
		return Vcs{}, fmt.Errorf("No repository found")
	}

	if metadata.IsFileSystemMetadataRoot(path) {
		return NewDefaultVcs(path)
	}

	return loadDefaultVcs(parent)
}

type Vcs struct {
	fs   metadata.Metadata
	root string
}

func (v Vcs) Init() error {
	if metadata.IsFileSystemMetadataRoot(v.root) {
		return fmt.Errorf("Repo already exists")
	}

	return engine.CreateNewIndex(v.fs)
}

func (v Vcs) Stat() (Status, error) {
	i, err := engine.LoadIndex(v.fs)
	if err != nil {
		return Status{}, err
	}

	return newStatusFromIndex(i), nil
}

func (v Vcs) AddFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	relPath, err := v.relativePath(path)
	if err != nil {
		return err
	}

	// TODO: Search commit history for existing file here.
	hash, err := engine.CreateNewDiff(v.fs, relPath, data)
	if err != nil {
		return err
	}

	return engine.AddDiffToIndex(v.fs, relPath, hash)
}

func (v Vcs) relativePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Rel(v.root, path)
	}

	return path, nil
}
