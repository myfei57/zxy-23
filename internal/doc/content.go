package doc

import (
	"fmt"

	"signflow/internal/model"
)

// writeContent durably replaces the document body of a contract revision.
// Callers must persist the revision marker only after this write succeeds.
func (s *Service) writeContent(contractID string, revision int, body string, now string) error {
	content := &model.DocumentContent{
		ContractID: contractID,
		Revision:   revision,
		Content:    body,
		UpdatedAt:  now,
	}
	if err := s.fs.WriteJSON(s.cfg.ContentFile(contractID), content); err != nil {
		return fmt.Errorf("doc: persist content %s: %w", contractID, err)
	}
	return nil
}
