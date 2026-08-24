package change

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

var lastDoc = map[string]*model.Contract{}

// View builds the change-history page model from the contract snapshot kept by
// the page session.
func (s *Service) View(contractID string) (*model.ChangeView, error) {
	contract := lastDoc[contractID]
	if contract == nil {
		current, err := s.doc.Current(contractID)
		if err != nil {
			return nil, err
		}
		contract = current
		lastDoc[contractID] = current
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
