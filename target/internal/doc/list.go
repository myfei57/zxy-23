package doc

import (
	"sort"
	"strings"

	"signflow/internal/model"
)

// List returns every contract of a namespace ordered by creation time.
func (s *Service) List(namespaceID string) ([]model.Contract, error) {
	names, err := s.fs.List("doc", "contracts")
	if err != nil {
		return nil, err
	}
	var contracts []model.Contract
	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		contractID := strings.TrimSuffix(name, ".json")
		contract, err := s.Current(contractID)
		if err != nil {
			continue
		}
		if contract.NamespaceID != namespaceID {
			continue
		}
		contracts = append(contracts, *contract)
	}
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].CreatedAt < contracts[j].CreatedAt
	})
	return contracts, nil
}

// Signers returns the signing parties of a contract in order.
func (s *Service) Signers(contractID string) ([]model.Signer, error) {
	contract, err := s.Current(contractID)
	if err != nil {
		return nil, err
	}
	return contract.Signers, nil
}
