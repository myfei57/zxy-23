package change

import (
	"fmt"

	"signflow/internal/model"
)

// Latest returns the most recent change journal row of a contract.
func (s *Service) Latest(contractID string) (*model.ChangeEntry, error) {
	var entries []model.ChangeEntry
	if err := s.fs.ReadJSON(s.cfg.ChangeJournalFile(contractID), &entries); err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("change: contract %s has no journal rows", contractID)
	}
	return &entries[len(entries)-1], nil
}
