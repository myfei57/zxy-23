package doc

import (
	"errors"
	"fmt"

	"signflow/internal/model"
)

// ErrNoChange reports a revision that does not alter the document body.
var ErrNoChange = errors.New("doc: revision does not change the document body")

// Revise advances a contract to the next revision. The new document body must
// be durable before the revision counter and the revised marker are persisted.
func (s *Service) Revise(contractID string, newContent string, now string) (*model.Contract, error) {
	contract, err := s.Current(contractID)
	if err != nil {
		return nil, err
	}
	if contract.Status == model.StatusArchived {
		return nil, fmt.Errorf("%w: cannot revise", ErrArchived)
	}
	if contract.Status == model.StatusEffective {
		return nil, fmt.Errorf("doc: effective contract cannot be revised")
	}
	current, err := s.Content(contractID)
	if err != nil {
		return nil, err
	}
	if current.Content == newContent {
		return nil, fmt.Errorf("%w: revision %d", ErrNoChange, contract.Revision)
	}
	nextRevision := contract.Revision + 1
	if err := s.writeContent(contractID, nextRevision, newContent, now); err != nil {
		return nil, err
	}
	contract.Revision = nextRevision
	contract.RevisedAt = now
	contract.UpdatedAt = now
	if err := s.persist(contract); err != nil {
		return nil, err
	}
	return contract, nil
}
