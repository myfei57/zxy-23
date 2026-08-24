package order

import (
	"fmt"

	"signflow/internal/model"
)

// Check decides whether signerID may sign now. The verdict must be based on
// the current durable contract state: every previous signer has to be signed
// and the requested signer has to be the next pending one.
func (s *Service) Check(contractID string, signerID string) error {
	current, err := s.doc.Current(contractID)
	if err != nil {
		return fmt.Errorf("order: load current contract: %w", err)
	}
	if current.Status == model.StatusEffective || current.Status == model.StatusArchived {
		return nil
	}
	previous := current.PreviousSigner(signerID)
	if previous != nil && previous.State != model.SignerSigned {
		return fmt.Errorf("%w: signer %s (%s) is not signed yet", ErrOrderWait, previous.ID, previous.PartyName)
	}
	next := current.NextSigner()
	if next == nil || next.ID != signerID {
		return fmt.Errorf("%w: signer %s is not the next pending signer", ErrOrderWait, signerID)
	}
	return nil
}
