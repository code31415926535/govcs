package engine

import (
	"testing"

	"github.com/code31415926535/govcs/metadata"

	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func createAndAddDiff(t *testing.T, fs metadata.Metadata, filename, content string) {
	diff, err := CreateNewDiff(fs, filename, []byte(content))
	assert.Nil(t, err, "could not create diff file")
	err = AddDiffToIndex(fs, filename, diff)
	assert.Nil(t, err, "could not add diff to index")
}

func createAndAddDiffAndCommit(t *testing.T, fs metadata.Metadata, filename, content, commitMsg string) {
	createAndAddDiff(t, fs, filename, content)
	commit, err := CreateNewCommit(fs, commitMsg)
	assert.Nil(t, err, "could not create commit file")
	err = ChangeHeadForce(fs, commit)
	assert.Nil(t, err, "could not clear index")
}

func Test_Commit_CanCommitOnceFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "could not create index file")
	createAndAddDiff(t, fs, "test.txt", "Hello World")
	commit, err := CreateNewCommit(fs, "initial commit")
	assert.Nil(t, err, "could not create commit file")

	commitFile := filepath.Join("commits", commit)
	assert.True(t, fs.FileExists(commitFile), "Commit is not created")
}

func Test_Commit_CanListAllCommits(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "could not create index file")
	createAndAddDiffAndCommit(t, fs, "test.txt", "Hello World", "initial commit")
	createAndAddDiffAndCommit(t, fs, "Readme.txt", "Dummy repo\n\nThis is a dummy repo", "adding readme")
	createAndAddDiffAndCommit(t, fs, "asdf.txt", "asdf", "adding asdf")

	commits, err := LoadCommits(fs, -1)
	assert.Nil(t, err, "could not load commits")
	assert.Equal(t, 3, len(commits), "could not load all commits")

	commits2, err := LoadCommits(fs, 2)
	assert.Nil(t, err, "could not load commits")
	assert.Equal(t, 2, len(commits2), "could not partially load commits")
}

func Test_Commit_FindLatestDiffOfFile(t *testing.T) {
	fs := metadata.NewInMemoryMetadata()

	err := CreateNewIndex(fs)
	assert.Nil(t, err, "could not create index file")
	createAndAddDiffAndCommit(t, fs, "test.txt", "Hello World", "initial commit")
	createAndAddDiffAndCommit(t, fs, "Readme.txt", "Dummy repo\n\nThis is a dummy repo", "adding readme")
	createAndAddDiffAndCommit(t, fs, "asdf.txt", "asdf", "adding asdf")

	diff, errDiff := FindLatestDiffOfFile(fs, "test.txt")
	assert.Nil(t, errDiff, "could not find latest diff")
	assert.NotEqual(t, "", diff, "no diff returned")

	diff2, errDiff2 := FindLatestDiffOfFile(fs, "none.txt")
	assert.Nil(t, errDiff2, "could not find latest diff")
	assert.Equal(t, "", diff2, "diff returned")
}
