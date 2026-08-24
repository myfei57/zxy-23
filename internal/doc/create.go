package doc

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/model"
)

// SignerInput is the creation payload for one signing party.
type SignerInput struct {
	PartyName string
	Email     string
}

// CreateInput is the contract creation payload used by the console API.
type CreateInput struct {
	NamespaceID string
	Title       string
	Signers     []SignerInput
	Content     string
	Now         string
}

// Create stores a new draft contract. The quota gate must run before any
// durable file is written so an over-quota contract never consumes storage.
func (s *Service) Create(input CreateInput) (*model.Contract, error) {
	if err := validateCreate(input); err != nil {
		return nil, err
	}
	contractID := uuid.NewString()
	contract := &model.Contract{
		ID:          contractID,
		NamespaceID: input.NamespaceID,
		Title:       input.Title,
		Status:      model.StatusDraft,
		Revision:    1,
		Signers:     buildSigners(contractID, input.Signers),
		CreatedAt:   input.Now,
		UpdatedAt:   input.Now,
	}
	if err := s.writeContent(contract.ID, 1, input.Content, input.Now); err != nil {
		return nil, err
	}
	if err := s.persist(contract); err != nil {
		return nil, err
	}
	if err := s.quota.Check(input.NamespaceID, 1); err != nil {
		return nil, err
	}
	if err := s.quota.Use(input.NamespaceID, 1, input.Now); err != nil {
		_ = s.fs.Remove(s.cfg.ContractFile(contract.ID)...)
		_ = s.fs.Remove(s.cfg.ContentFile(contract.ID)...)
		return nil, fmt.Errorf("doc: reserve quota after create: %w", err)
	}
	return contract, nil
}

func buildSigners(contractID string, inputs []SignerInput) []model.Signer {
	signers := make([]model.Signer, 0, len(inputs))
	for i, input := range inputs {
		signers = append(signers, model.Signer{
			ID:         uuid.NewString(),
			ContractID: contractID,
			PartyName:  input.PartyName,
			Email:      input.Email,
			Order:      i + 1,
			State:      model.SignerPending,
		})
	}
	return signers
}
