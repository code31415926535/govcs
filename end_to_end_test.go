package govcs_test

import (
	"crypto/md5"
	"fmt"
	"govcs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EndToEnd(t *testing.T) {
	// Setup
	repoRoot := createRepoRoot(t)

	t.Run("CanCreateRepoOnceAndOnlyOnce", func(t *testing.T) {
		repo, err := govcs.NewDefaultVcs(repoRoot)
		assert.Nil(t, err, "could not create repo")
		err2 := repo.Init()
		assert.Nil(t, err2, "could not create repo")

		err3 := repo.Init()
		assert.Equal(t, govcs.ErrRepoAlreadyExists, err3, "wrong error on second create")

		assert.Equal(t, "6a992d5529f459a44fee58c733255e86", calculateRepoHash(t, repoRoot))
	})

	t.Run("CanAddAndCommitNewFile", func(t *testing.T) {
		repo, err := govcs.LoadDefaultVcs(repoRoot)
		assert.Nil(t, err, "could not load repo")
		copyFile(t, "test.txt", "test.txt", repoRoot)

		errAdd := repo.AddFile(filepath.Join(repoRoot, "test.txt"))
		assert.Nil(t, errAdd, "could not add file")
		assertDiffCount(t, repoRoot, 1)
		errCommit := repo.CommitChanges("initial commit")
		assert.Nil(t, errCommit, "could not commit changes")
		assertDiffCount(t, repoRoot, 1)

		assert.Equal(t, "76ecfb4318c9581fab76f3295e6032f2", calculateRepoHash(t, repoRoot))
	})

	t.Run("CanAddThenRemoveChange", func(t *testing.T) {
		repo, err := govcs.LoadDefaultVcs(repoRoot)
		assert.Nil(t, err, "could not load repo")
		copyFile(t, "roses.poem", "roses.poem", repoRoot)

		errAdd := repo.AddFile(filepath.Join(repoRoot, "roses.poem"))
		assert.Nil(t, errAdd, "could not add file")
		assertDiffCount(t, repoRoot, 2)

		errRemove := repo.RemoveFile(filepath.Join(repoRoot, "roses.poem"))
		assert.Nil(t, errRemove, "could not remove file")
		assertDiffCount(t, repoRoot, 1)

		assert.Equal(t, "74efab7ddf22ce55ff90792fbed4f979", calculateRepoHash(t, repoRoot))
	})

	t.Run("CannotAddNonExistingFile", func(t *testing.T) {
		repo, err := govcs.LoadDefaultVcs(repoRoot)
		assert.Nil(t, err, "could not load repo")

		errAdd := repo.AddFile(filepath.Join(repoRoot, "none.txt"))
		assert.NotNil(t, errAdd, "could add non existent file")
	})

	t.Run("CannotRemoveNonExistingFile", func(t *testing.T) {
		repo, err := govcs.LoadDefaultVcs(repoRoot)
		assert.Nil(t, err, "could not load repo")

		errRemove := repo.RemoveFile(filepath.Join(repoRoot, "none.txt"))
		assert.NotNil(t, errRemove, "could add non existent file")

		assert.Equal(t, "74efab7ddf22ce55ff90792fbed4f979", calculateRepoHash(t, repoRoot))
	})

	// This test case is left here to help with debugging
	t.Run("ListResultingRepo", func(t *testing.T) {
		listDir(t, ".govcs", filepath.Join(repoRoot, ".govcs"))
		listDir(t, "diffs", filepath.Join(repoRoot, ".govcs", "diffs"))
		listDir(t, "commits", filepath.Join(repoRoot, ".govcs", "commits"))
	})

	// Teardown
	removeRepoRoot(t, repoRoot)
}

func copyFile(t *testing.T, src, dest, root string) {
	sourceFile := filepath.Join("test_data", src)
	destfile := filepath.Join(root, dest)

	input, err := ioutil.ReadFile(sourceFile)
	assert.Nil(t, err, "failed to copy file")
	err = ioutil.WriteFile(destfile, input, 0644)
	assert.Nil(t, err, "failed to copy file")
}

func createRepoRoot(t *testing.T) string {
	tmpDir := os.TempDir()
	dir, err := ioutil.TempDir(tmpDir, "govcs")
	assert.Nil(t, err, "Failed to create temp directory")
	return dir
}

func removeRepoRoot(t *testing.T, repoRoot string) {
	err := os.RemoveAll(repoRoot)
	assert.Nil(t, err, "Failed to remove perform cleanup")
}

func assertDiffCount(t *testing.T, repoRoot string, count int) {
	assertFileCountInDir(t, filepath.Join(repoRoot, ".govcs", "diffs"), count)
}

func assertCommitCount(t *testing.T, repoRoot string, count int) {
	assertFileCountInDir(t, filepath.Join(repoRoot, ".govcs", "commits"), count)
}

func assertFileCountInDir(t *testing.T, path string, count int) {
	files, err := ioutil.ReadDir(path)
	assert.Nil(t, err, "failed to read dir")
	assert.Equal(t, count, len(files), "count missmach")
}

func calculateRepoHash(t *testing.T, root string) string {
	return fmt.Sprintf("%x", md5.Sum(sumHierarchy(t, root)))
}

func sumHierarchy(t *testing.T, path string) []byte {
	stat, err := os.Stat(path)
	assert.Nil(t, err, "failed to sum hierarchy")
	if stat.IsDir() {
		files, err := ioutil.ReadDir(path)
		assert.Nil(t, err, "failed to sum hierarchy")

		var sum []byte
		for _, f := range files {
			sum = append(sum, sumHierarchy(t, filepath.Join(path, f.Name()))...)
		}
		return sum
	}

	return []byte(stat.Name())
}

func listDir(t *testing.T, prefix, path string) {
	files, err := ioutil.ReadDir(path)
	assert.Nil(t, err, "failed to read dir")

	t.Logf("%s\n", prefix)
	for _, f := range files {
		typeName := "-"
		if f.IsDir() {
			typeName = "d"
		}
		t.Logf("%s/%s\n", typeName, f.Name())
	}
	t.Log()
}
