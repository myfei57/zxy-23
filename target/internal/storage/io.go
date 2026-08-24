package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteJSON durably replaces the target file with the JSON encoding of value.
// The file is opened with truncation so a read-only target fails with a real
// write error instead of being silently replaced.
func (fs *FS) WriteJSON(parts []string, value any) error {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("storage: encode json: %w", err)
	}
	return fs.WriteFile(parts, payload)
}

// ReadJSON decodes the target file into value. A missing file yields a typed
// not-found error that callers can test with storage.IsNotExist.
func (fs *FS) ReadJSON(parts []string, value any) error {
	payload, err := fs.ReadFile(parts)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, value); err != nil {
		return fmt.Errorf("storage: decode json %s: %w", parts, err)
	}
	return nil
}

// WriteFile writes raw bytes to the target file, creating parent directories.
func (fs *FS) WriteFile(parts []string, data []byte) error {
	if err := fs.EnsureDir(parts...); err != nil {
		return err
	}
	path, err := fs.Path(parts...)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ReadFile returns the raw bytes of the target file.
func (fs *FS) ReadFile(parts []string) ([]byte, error) {
	path, err := fs.Path(parts...)
	if err != nil {
		return nil, err
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &NotFoundError{parts: parts}
		}
		return nil, err
	}
	return payload, nil
}

// AppendJSON reads the existing array at the target file, appends value and
// writes the file back. An empty or missing file starts a fresh array.
func (fs *FS) AppendJSON(parts []string, value any) error {
	var existing []json.RawMessage
	_ = fs.ReadJSON(parts, &existing)
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("storage: encode append value: %w", err)
	}
	existing = append(existing, payload)
	return fs.WriteJSON(parts, existing)
}
