package govcs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "govcs")
	assert.Nil(t, err, "Could not create temp dir")
	return dir
}

func Test_VCSCreateRepository_RepositoryDoesNotExist_CreatesEmptyRepository(t *testing.T) {
	dir := tempDir(t)
	err := NewRepository(dir)
	assert.Nil(t, err, "Create repository returns error")
	assert.DirExists(t, filepath.Join(dir, ".govcs"), "No metadata directory created")
}

func Test_VCSCreateRepository_RepositoryExists_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	err := NewRepository(dir)
	assert.Nil(t, err, "Create repository returns error")

	err2 := NewRepository(dir)
	assert.Error(t, err2, "Repository created when one already exists")
}
