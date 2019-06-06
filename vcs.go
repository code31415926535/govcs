package govcs

import (
	"fmt"
	"github.com/code31415926535/govcs/engine"
	"github.com/code31415926535/govcs/metadata"
	"io/ioutil"
	"path/filepath"
)

var (
	ErrRepoNotFound      = fmt.Errorf("No repo found")
	ErrRepoAlreadyExists = fmt.Errorf("Repo already exists")
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
		return Vcs{}, ErrRepoNotFound
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
		return ErrRepoAlreadyExists
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

func (v Vcs) RemoveFile(path string) error {
	relPath, err := v.relativePath(path)
	if err != nil {
		return err
	}

	return engine.RemoveDiffFromIndex(v.fs, relPath)
}

func (v Vcs) CommitChanges(message string) error {
	commit, err := engine.CreateNewCommit(v.fs, message)
	if err != nil {
		return err
	}

	return engine.ChangeHeadForce(v.fs, commit)
}

func (v Vcs) ListCommits() (Commits, error) {
	id, err := engine.LoadIndex(v.fs)
	if err != nil {
		return Commits{}, err
	}

	commits, err := engine.LoadCommits(v.fs, -1)
	if err != nil {
		return Commits{}, err
	}

	return newCommits(id.Head, commits), nil
}

func (v Vcs) relativePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Rel(v.root, path)
	}

	return path, nil
}
