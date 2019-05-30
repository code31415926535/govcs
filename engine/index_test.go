package engine

import (
	"govcs/metadata"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	emptyIndexFile   = `{"head":"","diffs":{}}`
	oneDiffIndexFile = `{"head":"","diffs":{"test.txt":"12345678"}}`
)

func assertIndexFile(t *testing.T, fs metadata.Metadata, expected string) {
	data, err := fs.ReadFile("index")
	assert.Nil(t, err, "Failed to read index file")
	assert.Equal(t, expected, string(data), "Index file mismatch")
}

func Test_Index_CanCreateNewIndexFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "Failed to create index file")
	assertIndexFile(t, fs, emptyIndexFile)
}

func Test_Index_CanAddDiffToIndexFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "Failed to create index file")
	err = AddDiffToIndex(fs, "test.txt", "12345678")
	assert.Nil(t, err, "Failed to add diff to index file")
	assertIndexFile(t, fs, oneDiffIndexFile)
}

func Test_Index_CanRemoveDiffFromIndexFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "Failed to create index file")
	err = AddDiffToIndex(fs, "test.txt", "12345678")
	assert.Nil(t, err, "Failed to add diff to index file")
	err = RemoveDiffFromIndex(fs, "test.txt")
	assert.Nil(t, err, "Failed to remove diff from index file")
	assertIndexFile(t, fs, emptyIndexFile)
}
