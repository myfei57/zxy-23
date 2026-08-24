package console

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"signflow/internal/archive"
	"signflow/internal/audit"
	"signflow/internal/change"
	"signflow/internal/doc"
	"signflow/internal/ns"
	"signflow/internal/order"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/sign"
	"signflow/internal/storage"
)

// Server wires every component behind the chi router and the embedded pages.
type Server struct {
	cfg     *settings.Settings
	fs      *storage.FS
	ns      *ns.Service
	quota   *quota.Service
	doc     *doc.Service
	order   *order.Service
	sign    *sign.Service
	archive *archive.Service
	change  *change.Service
	audit   *audit.Service
	router  *chi.Mux
}

// NewServer opens the file store and constructs every service dependency.
func NewServer(cfg *settings.Settings) (*Server, error) {
	fs, err := storage.Open(cfg.Root)
	if err != nil {
		return nil, fmt.Errorf("console: open storage: %w", err)
	}
	quotaSvc := quota.NewService(fs, cfg)
	nsSvc := ns.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	orderSvc := order.NewService(docSvc)
	auditSvc := audit.NewService(fs, cfg)
	signSvc := sign.NewService(fs, cfg, docSvc, orderSvc, auditSvc)
	archiveSvc := archive.NewService(fs, cfg, docSvc)
	changeSvc := change.NewService(fs, cfg, docSvc)
	server := &Server{
		cfg:     cfg,
		fs:      fs,
		ns:      nsSvc,
		quota:   quotaSvc,
		doc:     docSvc,
		order:   orderSvc,
		sign:    signSvc,
		archive: archiveSvc,
		change:  changeSvc,
		audit:   auditSvc,
	}
	server.router = server.routes()
	return server, nil
}

// Router returns the http.Handler served by the command entry point.
func (s *Server) Router() http.Handler {
	return s.router
}

// Close releases the server's resources; the file store needs no cleanup.
func (s *Server) Close() error {
	return nil
}
