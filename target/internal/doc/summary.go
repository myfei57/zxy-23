package doc

import "signflow/internal/model"

// Summary counts the contracts of a namespace by lifecycle status.
func (s *Service) Summary(namespaceID string) (model.ContractSummary, error) {
	contracts, err := s.List(namespaceID)
	if err != nil {
		return model.ContractSummary{}, err
	}
	summary := model.ContractSummary{NamespaceID: namespaceID, Total: len(contracts)}
	for i := range contracts {
		switch contracts[i].Status {
		case model.StatusDraft:
			summary.Draft++
		case model.StatusSigning:
			summary.Signing++
		case model.StatusEffective:
			summary.Effective++
		case model.StatusArchived:
			summary.Archived++
		}
	}
	return summary, nil
}
