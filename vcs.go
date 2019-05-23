package govcs

import (
	"os"
	"path/filepath"
)

func NewRepository(path string) error {
	repoRoot := filepath.Join(path, ".govcs")
	return os.Mkdir(repoRoot, 0644)
}
