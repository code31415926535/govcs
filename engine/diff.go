package engine

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/code31415926535/govcs/metadata"

	"github.com/gabstv/go-bsdiff/pkg/bsdiff"
	"github.com/gabstv/go-bsdiff/pkg/bspatch"
)

const diffDir = "diffs"

// CreateNewDiff creates a new diff object for a new file at path
// with data and save the newly created diff.
//
// The return value is the hash of the diff.
func CreateNewDiff(fs metadata.Metadata, path string, data []byte) (string, error) {
	d, err := newDiff(path, []byte{}, data, "")
	if err != nil {
		return "", err
	}

	return d.filename(), d.save(fs)
}

// CreateNewDiffOver creates a new diff object for an existing file at path
// with data. The prevDiff indicates the last revision of the file to be
// used when creating the diff. The diff object will be saved in fs metadata store.
//
// The return value is the hash of the diff.
func CreateNewDiffOver(fs metadata.Metadata, path, prevDiff string, data []byte) (string, error) {
	old, err := loadFileFromDiff(fs, prevDiff)
	if err != nil {
		return "", err
	}

	d, err := newDiff(path, old, data, prevDiff)
	if err != nil {
		return "", err
	}

	return d.filename(), d.save(fs)
}

// RemoveDiff removes diff object identified by hash. This only
// deletes the diff file, it does not check if the diff is
// referenced anywhere.
//
// FIXME: This can potentially cause issues where a diff can
// be deleted corrupting the repository.
func RemoveDiff(fs metadata.Metadata, hash string) error {
	return fs.RemoveFile(filepath.Join(diffDir, hash))
}

func loadFileFromDiff(fs metadata.Metadata, hash string) ([]byte, error) {
	d, err := loadDiff(fs, hash)
	if err != nil {
		return nil, err
	}

	if d.PrevDiff == "" {
		return d.apply(nil)
	}

	base, err := loadFileFromDiff(fs, d.PrevDiff)
	if err != nil {
		return nil, err
	}

	return d.apply(base)
}

func loadDiff(fs metadata.Metadata, hash string) (*diff, error) {
	d, err := fs.ReadFile(filepath.Join(diffDir, hash))
	if err != nil {
		return nil, err
	}
	// Split
	parts := bytes.Split(d, []byte("\000BSDIFF40"))
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid diff: %s", hash)
	}

	meta := bytes.Split(parts[0], []byte("\000"))
	if len(meta) != 2 {
		return nil, fmt.Errorf("invalid diff: %s", hash)
	}

	return &diff{
		Path:     string(meta[0]),
		PrevDiff: string(meta[1]),
		Data:     append([]byte("BSDIFF40"), parts[1]...),
	}, nil
}

func newDiff(path string, old, new []byte, prevDiff string) (*diff, error) {
	b, err := bsdiff.Bytes(old, new)
	return &diff{path, prevDiff, b}, err
}

type diff struct {
	Path     string
	PrevDiff string
	Data     []byte
}

func (d *diff) apply(data []byte) ([]byte, error) {
	return bspatch.Bytes(data, d.Data)
}

func (d *diff) save(fs metadata.Metadata) error {
	return fs.WriteFile(filepath.Join(diffDir, d.filename()), d.data())
}

func (d *diff) filename() string {
	data := append([]byte(d.Path), []byte(d.PrevDiff)...)
	data = append(data, d.Data...)
	return dataHash(data)
}

func (d *diff) data() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s\000%s\000", d.Path, d.PrevDiff)
	b.Write(d.Data)
	return b.Bytes()
}
