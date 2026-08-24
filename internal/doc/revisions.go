package doc

import (
	"signflow/internal/model"
)

// RevisionState describes the durable pairing of revision counter and marker.
type RevisionState struct {
	ContractID string               `json:"contract_id"`
	Revision   int                  `json:"revision"`
	RevisedAt  string               `json:"revised_at"`
	Status     model.ContractStatus `json:"status"`
}

// RevisionState returns the current revision pairing of a contract; the
// change-history page uses it to explain what generation a document is on.
func (s *Service) RevisionState(contractID string) (RevisionState, error) {
	contract, err := s.Current(contractID)
	if err != nil {
		return RevisionState{}, err
	}
	return RevisionState{
		ContractID: contractID,
		Revision:   contract.Revision,
		RevisedAt:  contract.RevisedAt,
		Status:     contract.Status,
	}, nil
}
