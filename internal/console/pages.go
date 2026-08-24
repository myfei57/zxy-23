package console

import (
	"net/http"

	"signflow/internal/web"
)

func (s *Server) handleContractsPage(w http.ResponseWriter, r *http.Request) {
	s.servePage(w, "contracts.html")
}

func (s *Server) handleSigningPage(w http.ResponseWriter, r *http.Request) {
	s.servePage(w, "signing.html")
}

func (s *Server) handleArchivePage(w http.ResponseWriter, r *http.Request) {
	s.servePage(w, "archive.html")
}

func (s *Server) handleAuditPage(w http.ResponseWriter, r *http.Request) {
	s.servePage(w, "audit.html")
}

func (s *Server) servePage(w http.ResponseWriter, name string) {
	payload, err := web.Pages.ReadFile(name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "page not found")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(payload)
}
