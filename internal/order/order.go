package order

import (
	"errors"
	"fmt"

	"signflow/internal/doc"
)

// Service enforces the signing sequence of a contract.
type Service struct {
	doc *doc.Service
}

// NewService builds the order service on top of the document service.
func NewService(docSvc *doc.Service) *Service {
	return &Service{doc: docSvc}
}

// ErrOrderWait is returned when the requested signer is not allowed to act yet.
var ErrOrderWait = errors.New("order: previous signer has not finished")

// Open validates that a signing session can target the contract.
func (s *Service) Open(contractID string) error {
	if _, err := s.doc.Current(contractID); err != nil {
		return fmt.Errorf("order: open session: %w", err)
	}
	return nil
}
