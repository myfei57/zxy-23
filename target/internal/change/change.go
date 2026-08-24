package change

import (
	"signflow/internal/doc"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// Service owns the immutable change journal and the change-history view.
type Service struct {
	fs  *storage.FS
	cfg *settings.Settings
	doc *doc.Service
}

// NewService builds the change service over the document service.
func NewService(fs *storage.FS, cfg *settings.Settings, docSvc *doc.Service) *Service {
	return &Service{fs: fs, cfg: cfg, doc: docSvc}
}
