package storage

import (
	"os"
	"path/filepath"
)

func writeLocal(dir, key string, body []byte) error {
	full := filepath.Join(dir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, body, 0o644)
}

func readLocal(dir, key string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, filepath.FromSlash(key)))
}
