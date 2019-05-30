package metadata

import "fmt"

type inMemoryMetadata struct {
	files map[string][]byte
}

func (meta *inMemoryMetadata) ReadFile(path string) ([]byte, error) {
	if e, ok := meta.files[path]; ok {
		return e, nil
	}

	return nil, fmt.Errorf("File not found")
}

func (meta *inMemoryMetadata) WriteFile(path string, data []byte) error {
	meta.files[path] = data
	return nil
}

func (meta *inMemoryMetadata) FileExists(path string) bool {
	_, ok := meta.files[path]
	return ok
}
