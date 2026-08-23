package console

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(s.logRequests)
	router.Use(middleware.Recoverer)
	router.Get("/health", s.handleHealth)
	router.Get("/", s.handleContractsPage)
	router.Get("/contracts", s.handleContractsPage)
	router.Get("/signing", s.handleSigningPage)
	router.Get("/archive", s.handleArchivePage)
	router.Get("/audit", s.handleAuditPage)
	router.Route("/api", func(api chi.Router) {
		api.Post("/namespaces", s.handleCreateNamespace)
		api.Get("/namespaces", s.handleListNamespaces)
		api.Get("/stats", s.handleStats)
		api.Get("/namespaces/{id}/quota", s.handleNamespaceQuota)
		api.Get("/namespaces/{id}/dashboard", s.handleNamespaceDashboard)
		api.Post("/contracts", s.handleCreateContract)
		api.Get("/contracts", s.handleListContracts)
		api.Get("/archive/pending", s.handleArchivePending)
		api.Get("/contracts/{id}", s.handleGetContract)
		api.Get("/contracts/{id}/content", s.handleContractContent)
		api.Get("/contracts/{id}/signers", s.handleContractSigners)
		api.Get("/contracts/{id}/signing", s.handleSigningSummary)
		api.Get("/contracts/{id}/order", s.handleContractOrder)
		api.Get("/contracts/{id}/order/{signer_id}", s.handleSignerPosition)
		api.Get("/contracts/{id}/revision", s.handleContractRevision)
		api.Get("/contracts/{id}/records", s.handleContractRecords)
		api.Get("/contracts/{id}/records/latest", s.handleContractLastRecord)
		api.Get("/contracts/{id}/acks", s.handleContractAcks)
		api.Get("/contracts/{id}/cursor", s.handleContractCursor)
		api.Get("/contracts/{id}/changes", s.handleContractChanges)
		api.Get("/contracts/{id}/changes/latest", s.handleLatestChange)
		api.Get("/contracts/{id}/audit", s.handleContractAudit)
		api.Get("/contracts/{id}/audit/recent", s.handleContractAuditRecent)
		api.Get("/contracts/{id}/audit/latest", s.handleLatestAudit)
		api.Get("/contracts/{id}/timeline", s.handleContractTimeline)
		api.Get("/contracts/{id}/archive", s.handleArchiveDetail)
		api.Post("/contracts/{id}/archive/verify", s.handleArchiveVerify)
		api.Post("/contracts/{id}/revise", s.handleReviseContract)
		api.Post("/contracts/{id}/sign", s.handleSignContract)
		api.Post("/contracts/{id}/notify", s.handleNotifySigners)
		api.Post("/contracts/{id}/archive", s.handleArchiveContract)
		api.Post("/contracts/{id}/archive/retry", s.handleArchiveRetry)
		api.Get("/archive/records", s.handleArchiveRecords)
	})
	return router
}
