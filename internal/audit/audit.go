package audit

import (
	"signflow/internal/model"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// Service appends immutable audit events per entity.
type Service struct {
	fs  *storage.FS
	cfg *settings.Settings
}

// NewService builds the audit service over the shared file store.
func NewService(fs *storage.FS, cfg *settings.Settings) *Service {
	return &Service{fs: fs, cfg: cfg}
}

// List returns all audit events of one entity in append order.
func (s *Service) List(entityID string) ([]model.AuditEvent, error) {
	var events []model.AuditEvent
	if err := s.fs.ReadJSON(s.cfg.AuditEventFile(entityID), &events); err != nil {
		if storage.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return events, nil
}
