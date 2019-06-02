package engine

import (
	"bytes"
	"govcs/metadata"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Diff_CanCreateNewDiffForFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	_, err := CreateNewDiff(fs, "test.txt", []byte("Hello World"))
	assert.Nil(t, err, "Failed to create new diff")
	expectedFilepath := filepath.Join("diffs", "b600b7b768d5f376efa0aa8320b13245")
	assert.True(t, fs.FileExists(expectedFilepath), "Diff file not created")
	d, err := fs.ReadFile(expectedFilepath)
	assert.Nil(t, err, "Failed to read diff file")
	parts := bytes.Split(d, []byte("\000BSDIFF40"))
	partOne := bytes.Split(parts[0], []byte("\000"))
	assert.Equal(t, "test.txt", string(partOne[0]), "wrong diff")
	assert.Equal(t, 0, len(partOne[1]), "wrong diff")
}

func Test_Diff_CanCreateNewDiffOverExistingDiff(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	diff1, err := CreateNewDiff(fs, "test.txt", []byte("Hello World"))
	assert.Nil(t, err, "Failed to create new diff")

	diff2, err := CreateNewDiffOver(fs, "test.txt", diff1, []byte("Hello Universe"))
	assert.Nil(t, err, "Failed to create new diff over existing one")

	diff3, err := CreateNewDiffOver(fs, "test.txt", diff2, []byte("Hello Universe\nHow are you doing today?"))
	assert.Nil(t, err, "Failed to create new diff over existing one")

	d, err := fs.ReadFile(filepath.Join("diffs", diff3))
	assert.Nil(t, err, "Failed to read diff file")
	parts := bytes.Split(d, []byte("\000BSDIFF40"))
	partOne := bytes.Split(parts[0], []byte("\000"))
	assert.Equal(t, "test.txt", string(partOne[0]), "wrong diff")
	assert.Equal(t, diff2, string(partOne[1]), "wrong diff")
}
