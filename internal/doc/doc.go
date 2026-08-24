package doc

import (
	"errors"
	"fmt"

	"signflow/internal/model"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// ErrArchived reports a write attempted against an archived contract.
var ErrArchived = errors.New("doc: contract is archived")

// ErrSignerUnknown reports a signer that is not part of the contract.
var ErrSignerUnknown = errors.New("doc: signer is not part of the contract")

// ErrNotEffective reports an archive attempt against a non-effective contract.
var ErrNotEffective = errors.New("doc: contract is not effective")

// Service owns the contract aggregate, its durable content and its signing
// markers. Persistence is file-backed through storage.FS.
type Service struct {
	fs    *storage.FS
	cfg   *settings.Settings
	quota *quota.Service
}

// NewService builds the document service with the quota dependency used by
// the contract creation gate.
func NewService(fs *storage.FS, cfg *settings.Settings, quotaSvc *quota.Service) *Service {
	return &Service{fs: fs, cfg: cfg, quota: quotaSvc}
}

func (s *Service) persist(contract *model.Contract) error {
	if err := s.fs.WriteJSON(s.cfg.ContractFile(contract.ID), contract); err != nil {
		return fmt.Errorf("doc: persist contract %s: %w", contract.ID, err)
	}
	return nil
}

func (s *Service) loadContract(contractID string) (*model.Contract, error) {
	var contract model.Contract
	if err := s.fs.ReadJSON(s.cfg.ContractFile(contractID), &contract); err != nil {
		return nil, fmt.Errorf("doc: load contract %s: %w", contractID, err)
	}
	return &contract, nil
}
