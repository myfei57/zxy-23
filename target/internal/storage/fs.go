package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FS is the single persistence primitive of the platform. Every component
// writes JSON aggregates and blobs through this type; the physical layout is
// supplied by settings.Settings.
type FS struct {
	root string
}

// Open creates the root directory when needed and returns a handle bound to it.
func Open(root string) (*FS, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("storage: create root %s: %w", root, err)
	}
	return &FS{root: root}, nil
}

// Path resolves a relative file path and verifies it stays inside the root.
func (fs *FS) Path(parts ...string) (string, error) {
	all := append([]string{fs.root}, parts...)
	joined := filepath.Join(all...)
	clean := filepath.Clean(joined)
	rootClean := filepath.Clean(fs.root)
	if clean != rootClean && !strings.HasPrefix(clean, rootClean+string(os.PathSeparator)) {
		return "", fmt.Errorf("storage: path escapes root: %s", clean)
	}
	return clean, nil
}

// Exists reports whether the target file exists on disk.
func (fs *FS) Exists(parts ...string) bool {
	path, err := fs.Path(parts...)
	if err != nil {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// EnsureDir creates every directory component leading to the target file.
func (fs *FS) EnsureDir(parts ...string) error {
	path, err := fs.Path(parts...)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

// Remove deletes the target file when present.
func (fs *FS) Remove(parts ...string) error {
	path, err := fs.Path(parts...)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// List returns the names of regular files directly below a relative directory.
func (fs *FS) List(parts ...string) ([]string, error) {
	path, err := fs.Path(parts...)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	return names, nil
}
