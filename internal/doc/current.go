package doc

import (
	"fmt"

	"signflow/internal/model"
)

// Current returns the latest durable contract aggregate.
func (s *Service) Current(contractID string) (*model.Contract, error) {
	return s.loadContract(contractID)
}

// Content returns the latest durable document body of a contract.
func (s *Service) Content(contractID string) (*model.DocumentContent, error) {
	var content model.DocumentContent
	if err := s.fs.ReadJSON(s.cfg.ContentFile(contractID), &content); err != nil {
		return nil, fmt.Errorf("doc: load content %s: %w", contractID, err)
	}
	return &content, nil
}
