package engine

import (
	"govcs/metadata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Commit_CanCommitOnceFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "could not create index file")
	diff, err := CreateNewDiff(fs, "test.txt", []byte("Hello World"))
	assert.Nil(t, err, "could not create diff file")
	err = AddDiffToIndex(fs, "test.txt", diff)
	assert.Nil(t, err, "could not add diff to index")
	_, err = CreateNewCommit(fs, "initial commit")
	assert.Nil(t, err, "could not create commit file")

	// TODO: properly assert here
}
