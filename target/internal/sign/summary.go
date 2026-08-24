package sign

import (
	"signflow/internal/model"
)

// Summary builds the signing-page read model for one contract.
func (s *Service) Summary(contractID string) (model.SigningSummary, error) {
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return model.SigningSummary{}, err
	}
	signed := 0
	for i := range contract.Signers {
		if contract.Signers[i].State == model.SignerSigned {
			signed++
		}
	}
	return model.SigningSummary{
		ContractID: contractID,
		Total:      len(contract.Signers),
		Signed:     signed,
		Pending:    len(contract.Signers) - signed,
		NextSigner: contract.NextSigner(),
		Complete:   contract.AllSigned(),
	}, nil
}
