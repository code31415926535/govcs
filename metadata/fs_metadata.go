package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const defaultPerm = 0644

type actualMetadata struct {
	root string
}

func (meta *actualMetadata) ReadFile(path string) ([]byte, error) {
	absPath := filepath.Join(meta.root, path)
	return ioutil.ReadFile(absPath)
}

func (meta *actualMetadata) WriteFile(path string, data []byte) error {
	absPath := filepath.Join(meta.root, path)
	err := os.MkdirAll(filepath.Dir(absPath), defaultPerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(absPath, data, defaultPerm)
}

func (meta *actualMetadata) FileExists(path string) bool {
	absPath := filepath.Join(meta.root, path)
	if info, err := os.Stat(absPath); err == nil && info.Mode().IsRegular() {
		return true
	}

	return false
}
