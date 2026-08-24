package sign

import "signflow/internal/model"

// Signers returns the signing parties of a contract in order.
func (s *Service) Signers(contractID string) ([]model.Signer, error) {
	return s.doc.Signers(contractID)
}
