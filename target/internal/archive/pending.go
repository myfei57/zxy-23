package archive

import "signflow/internal/model"

// Pending returns the effective contracts of a namespace that still lack an
// archive ledger row; the nightly archive batch starts from this worklist.
func (s *Service) Pending(namespaceID string) ([]model.Contract, error) {
	contracts, err := s.doc.List(namespaceID)
	if err != nil {
		return nil, err
	}
	pending := make([]model.Contract, 0, len(contracts))
	for _, contract := range contracts {
		if contract.Status != model.StatusEffective {
			continue
		}
		recorded, err := s.hasRecord(contract.ID)
		if err != nil {
			return nil, err
		}
		if !recorded {
			pending = append(pending, contract)
		}
	}
	return pending, nil
}
