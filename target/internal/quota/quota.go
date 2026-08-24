package quota

import (
	"fmt"

	"signflow/internal/model"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// Service owns the durable capacity ledger of each namespace.
type Service struct {
	fs  *storage.FS
	cfg *settings.Settings
}

// NewService builds the quota service over the shared file store.
func NewService(fs *storage.FS, cfg *settings.Settings) *Service {
	return &Service{fs: fs, cfg: cfg}
}

// Ensure creates the ledger for a namespace when it is missing.
func (s *Service) Ensure(namespaceID string, limit int, now string) error {
	if s.fs.Exists(s.cfg.QuotaLedgerFile(namespaceID)...) {
		return nil
	}
	ledger := &model.QuotaLedger{
		NamespaceID: namespaceID,
		Limit:       limit,
		Used:        0,
		UpdatedAt:   now,
	}
	if err := s.fs.WriteJSON(s.cfg.QuotaLedgerFile(namespaceID), ledger); err != nil {
		return fmt.Errorf("quota: persist ledger: %w", err)
	}
	return nil
}

// Use reserves delta capacity and persists the new used total.
func (s *Service) Use(namespaceID string, delta int, now string) error {
	ledger, err := s.load(namespaceID)
	if err != nil {
		return err
	}
	ledger.Used += delta
	ledger.UpdatedAt = now
	if err := s.fs.WriteJSON(s.cfg.QuotaLedgerFile(namespaceID), ledger); err != nil {
		return fmt.Errorf("quota: persist usage: %w", err)
	}
	return nil
}

// Snapshot returns the API read model of one namespace ledger.
func (s *Service) Snapshot(namespaceID string) (model.QuotaSnapshot, error) {
	ledger, err := s.load(namespaceID)
	if err != nil {
		return model.QuotaSnapshot{}, err
	}
	return ledger.Snapshot(), nil
}

func (s *Service) load(namespaceID string) (*model.QuotaLedger, error) {
	var ledger model.QuotaLedger
	if err := s.fs.ReadJSON(s.cfg.QuotaLedgerFile(namespaceID), &ledger); err != nil {
		return nil, fmt.Errorf("quota: load ledger %s: %w", namespaceID, err)
	}
	return &ledger, nil
}
