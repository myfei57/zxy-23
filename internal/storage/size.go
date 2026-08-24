package storage

import (
	"os"
	"time"
)

// Size returns the byte size of a target file.
func (fs *FS) Size(parts ...string) (int64, error) {
	path, err := fs.Path(parts...)
	if err != nil {
		return 0, err
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNotExist
		}
		return 0, err
	}
	return info.Size(), nil
}

// ModTime returns the modification time of a target file.
func (fs *FS) ModTime(parts ...string) (time.Time, error) {
	path, err := fs.Path(parts...)
	if err != nil {
		return time.Time{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, ErrNotExist
		}
		return time.Time{}, err
	}
	return info.ModTime().UTC(), nil
}

// Describe joins the size and modification time into a diagnostics row used by
// the archive verification read model.
func (fs *FS) Describe(parts ...string) (int64, time.Time, error) {
	size, err := fs.Size(parts...)
	if err != nil {
		return 0, time.Time{}, err
	}
	modified, err := fs.ModTime(parts...)
	if err != nil {
		return 0, time.Time{}, err
	}
	return size, modified, nil
}
