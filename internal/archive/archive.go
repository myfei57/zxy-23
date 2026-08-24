package archive

import (
	"signflow/internal/doc"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// Service produces durable archive artifacts and the archive ledger.
type Service struct {
	fs  *storage.FS
	cfg *settings.Settings
	doc *doc.Service
}

// NewService builds the archive service over the document service.
func NewService(fs *storage.FS, cfg *settings.Settings, docSvc *doc.Service) *Service {
	return &Service{fs: fs, cfg: cfg, doc: docSvc}
}
