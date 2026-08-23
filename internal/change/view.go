package change

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

// View builds the change-history page model. The current revision must be read
// from the live document store, never from a snapshot captured earlier.
func (s *Service) View(contractID string) (*model.ChangeView, error) {
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	var entries []model.ChangeEntry
	if err := s.fs.ReadJSON(s.cfg.ChangeJournalFile(contractID), &entries); err != nil {
		if storage.IsNotFound(err) {
			entries = nil
		} else {
			return nil, err
		}
	}
	return &model.ChangeView{
		ContractID:      contractID,
		CurrentRevision: contract.Revision,
		Title:           contract.Title,
		Entries:         entries,
	}, nil
}
