package metadata

import (
	"os"
	"path/filepath"
)

const (
	metadataDir = ".govcs"
)

func NewFileSystemMetadata(root string) Metadata {
	return &actualMetadata{
		root: filepath.Join(root, metadataDir),
	}
}

func IsFileSystemMetadataRoot(path string) bool {
	m := filepath.Join(path, metadataDir)
	if info, err := os.Stat(m); err == nil && info.IsDir() {
		return true
	}

	return false
}

func NewInMemoryMetadata() Metadata {
	return &inMemoryMetadata{
		files: make(map[string][]byte),
	}
}

// Metadata is an abstractization of the filesystem that
//  that stores the metadata for the repo. This is used
//	to allow for easier testing (using in memory metadata)
//	while allowing the use of file system for production.
//	All paths have to be relative.
type Metadata interface {
	// Readfile reads the file located at path and
	//	returns it's content.
	ReadFile(path string) ([]byte, error)

	// Writefile writes file at path with data.
	//	Should make sure path exists and file exists.
	//	If file already exists, it should be overwritten.
	WriteFile(path string, data []byte) error

	// RemoveFile removes file at path.
	RemoveFile(path string) error

	// Returns if the current file exists or not
	FileExists(path string) bool
}
