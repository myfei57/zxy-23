package doc

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidInput reports a creation payload that cannot be stored.
var ErrInvalidInput = errors.New("doc: invalid contract input")

// validateCreate checks the creation payload before any file is touched.
func validateCreate(input CreateInput) error {
	if strings.TrimSpace(input.NamespaceID) == "" {
		return fmt.Errorf("%w: namespace is required", ErrInvalidInput)
	}
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if strings.TrimSpace(input.Content) == "" {
		return fmt.Errorf("%w: content is required", ErrInvalidInput)
	}
	if len(input.Signers) == 0 {
		return fmt.Errorf("%w: at least one signer is required", ErrInvalidInput)
	}
	seen := make(map[string]bool, len(input.Signers))
	for _, signer := range input.Signers {
		name := strings.TrimSpace(signer.PartyName)
		email := strings.TrimSpace(signer.Email)
		if name == "" {
			return fmt.Errorf("%w: signer name is required", ErrInvalidInput)
		}
		if email == "" {
			return fmt.Errorf("%w: signer email is required", ErrInvalidInput)
		}
		if seen[email] {
			return fmt.Errorf("%w: duplicate signer email %s", ErrInvalidInput, email)
		}
		seen[email] = true
	}
	return nil
}
