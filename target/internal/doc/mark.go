package doc

import (
	"fmt"

	"signflow/internal/model"
)

// MarkSigned persists the signed state of one signer on the contract and
// advances the contract to effective once every signer has signed.
func (s *Service) MarkSigned(contractID string, signerID string, cert string, signTime string, at string) error {
	contract, err := s.Current(contractID)
	if err != nil {
		return err
	}
	if contract.Status == model.StatusArchived {
		return fmt.Errorf("%w: cannot sign", ErrArchived)
	}
	if contract.Status == model.StatusEffective {
		return nil
	}
	index := contract.SignerIndex(signerID)
	if index < 0 {
		return fmt.Errorf("%w: %s", ErrSignerUnknown, signerID)
	}
	signer := &contract.Signers[index]
	signer.MarkSigned(cert, signTime, at)
	if contract.Status == model.StatusDraft {
		contract.Status = model.StatusSigning
	}
	if contract.AllSigned() {
		contract.Status = model.StatusEffective
		contract.EffectiveAt = at
		contract.SignedAt = at
	}
	contract.UpdatedAt = at
	return s.persist(contract)
}

// MarkArchived persists the archived marker of a contract.
func (s *Service) MarkArchived(contractID string, at string) error {
	contract, err := s.Current(contractID)
	if err != nil {
		return err
	}
	if contract.Status != model.StatusEffective {
		return fmt.Errorf("%w: %s", ErrNotEffective, contract.Status)
	}
	contract.Status = model.StatusArchived
	contract.ArchivedAt = at
	contract.UpdatedAt = at
	return s.persist(contract)
}
