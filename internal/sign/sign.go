package sign

import (
	"errors"

	"signflow/internal/audit"
	"signflow/internal/doc"
	"signflow/internal/order"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// ErrClosed reports a signing attempt against a finished contract.
var ErrClosed = errors.New("sign: contract is closed for signing")

// ErrMissingCertificate reports a signature without certificate evidence.
var ErrMissingCertificate = errors.New("sign: certificate fingerprint is required")

// Service orchestrates signature records, signed markers and invitation
// acknowledgements for one contract.
type Service struct {
	fs    *storage.FS
	cfg   *settings.Settings
	doc   *doc.Service
	order *order.Service
	audit *audit.Service
}

// NewService builds the signing service with its document, order and audit
// dependencies.
func NewService(fs *storage.FS, cfg *settings.Settings, docSvc *doc.Service, orderSvc *order.Service, auditSvc *audit.Service) *Service {
	return &Service{fs: fs, cfg: cfg, doc: docSvc, order: orderSvc, audit: auditSvc}
}
