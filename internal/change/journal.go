package change

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"signflow/internal/model"
)

// Append records one immutable change row for a contract.
func (s *Service) Append(contractID string, action string, revision int, prevRevision int, note string, now string) error {
	if strings.TrimSpace(action) == "" {
		return fmt.Errorf("change: action is required")
	}
	entry := &model.ChangeEntry{
		ID:           uuid.NewString(),
		ContractID:   contractID,
		Action:       action,
		Revision:     revision,
		PrevRevision: prevRevision,
		Note:         note,
		ChangedAt:    now,
	}
	if err := s.fs.AppendJSON(s.cfg.ChangeJournalFile(contractID), entry); err != nil {
		return fmt.Errorf("change: persist journal row: %w", err)
	}
	return nil
}
