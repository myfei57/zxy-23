package doc

import (
	"strings"

	"signflow/internal/model"
)

// Search returns the contracts of a namespace whose title or any signer
// matches the keyword. An empty keyword returns the full list.
func (s *Service) Search(namespaceID string, keyword string) ([]model.Contract, error) {
	contracts, err := s.List(namespaceID)
	if err != nil {
		return nil, err
	}
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	if keyword == "" {
		return contracts, nil
	}
	matches := make([]model.Contract, 0, len(contracts))
	for _, contract := range contracts {
		if strings.Contains(strings.ToLower(contract.Title), keyword) {
			matches = append(matches, contract)
			continue
		}
		for _, signer := range contract.Signers {
			if strings.Contains(strings.ToLower(signer.PartyName), keyword) ||
				strings.Contains(strings.ToLower(signer.Email), keyword) {
				matches = append(matches, contract)
				break
			}
		}
	}
	return matches, nil
}
