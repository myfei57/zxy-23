package order

import (
	"errors"
	"fmt"

	"signflow/internal/model"
)

// ErrSignerNotInOrder reports a signer that does not belong to the contract.
var ErrSignerNotInOrder = errors.New("order: signer is not part of the contract")

// Sequence returns the signing order read model of a contract, including which
// signer is next and who is blocking the queue.
func (s *Service) Sequence(contractID string) ([]model.OrderPosition, error) {
	current, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	positions := make([]model.OrderPosition, 0, len(current.Signers))
	for index := range current.Signers {
		signer := current.Signers[index]
		position := model.OrderPosition{
			SignerID:  signer.ID,
			PartyName: signer.PartyName,
			Position:  signer.Order,
			Total:     len(current.Signers),
			State:     signer.State,
		}
		position.IsNext = position.Position == 1 && signer.State == model.SignerPending
		if index > 0 && current.Signers[index-1].State != model.SignerSigned {
			position.BlockedBy = current.Signers[index-1].PartyName
		}
		positions = append(positions, position)
	}
	return positions, nil
}

// Position returns one signer's place in the signing sequence.
func (s *Service) Position(contractID string, signerID string) (model.OrderPosition, error) {
	positions, err := s.Sequence(contractID)
	if err != nil {
		return model.OrderPosition{}, err
	}
	for _, position := range positions {
		if position.SignerID == signerID {
			return position, nil
		}
	}
	return model.OrderPosition{}, fmt.Errorf("%w: %s", ErrSignerNotInOrder, signerID)
}
