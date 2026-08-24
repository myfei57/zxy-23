package storage

import (
	"errors"
	"fmt"
	"strings"
)

// NotFoundError reports a missing persistence target.
type NotFoundError struct {
	parts []string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("storage: record not found: %s", strings.Join(e.parts, "/"))
}

func (e *NotFoundError) Unwrap() error {
	return ErrNotExist
}

// IsNotFound reports whether err is a storage not-found error.
func IsNotFound(err error) bool {
	var target *NotFoundError
	return errors.As(err, &target)
}

// ErrNotExist is returned by convenience wrappers when a target is absent.
var ErrNotExist = errors.New("storage: target does not exist")
