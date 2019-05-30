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
	assert.Equal(t, 9, len(parts[0]), "wrong split")
}

func Test_Diff_CanCreateNewDiffOverExistingDiff(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	_, err := CreateNewDiff(fs, "test.txt", []byte("Hello World"))
	assert.Nil(t, err, "Failed to create new diff")

	err = CreateNewDiffOver(fs, "test.txt", "b600b7b768d5f376efa0aa8320b13245", []byte("Hello Universe"))
	assert.Nil(t, err, "Failed to create new diff over existing one")

	err = CreateNewDiffOver(fs, "test.txt", "a8dee22bf2fb8bffd56b1c33fd043c80", []byte("Hello Universe\nHow are you doing today?"))
	assert.Nil(t, err, "Failed to create new diff over existing one")
}
