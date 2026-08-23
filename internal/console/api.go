package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"signflow/internal/doc"
	"signflow/internal/model"
	"signflow/internal/order"
	"signflow/internal/storage"
)

func now() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func readJSON(r *http.Request, value any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("console: decode request: %w", err)
	}
	return nil
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case storage.IsNotFound(err):
		writeError(w, http.StatusNotFound, "record not found")
	case errors.Is(err, order.ErrOrderWait):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, doc.ErrArchived):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, doc.ErrSignerUnknown):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, doc.ErrNotEffective):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "signflow"})
}

func (s *Server) handleCreateNamespace(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		LegalEntity string `json:"legal_entity"`
		QuotaLimit  int    `json:"quota_limit"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.QuotaLimit <= 0 {
		input.QuotaLimit = 100
	}
	namespace, err := s.ns.Create(input.Name, input.LegalEntity, now())
	if err != nil {
		handleError(w, err)
		return
	}
	if err := s.quota.Ensure(namespace.ID, input.QuotaLimit, now()); err != nil {
		handleError(w, err)
		return
	}
	_, _ = s.audit.Record(namespace.ID, "namespace.created", "success", input.Name, now())
	writeJSON(w, http.StatusCreated, namespace)
}

func (s *Server) handleListNamespaces(w http.ResponseWriter, r *http.Request) {
	namespaces, err := s.ns.List()
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, namespaces)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	count, err := s.ns.Count()
	if err != nil {
		handleError(w, err)
		return
	}
	auditTotal, err := s.audit.Total()
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"namespaces": count, "audit_events": auditTotal})
}

func (s *Server) handleNamespaceQuota(w http.ResponseWriter, r *http.Request) {
	namespaceID := chi.URLParam(r, "id")
	snapshot, err := s.quota.Snapshot(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) handleNamespaceDashboard(w http.ResponseWriter, r *http.Request) {
	namespaceID := chi.URLParam(r, "id")
	namespace, err := s.ns.Get(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	quotaSnapshot, err := s.quota.Snapshot(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	summary, err := s.doc.Summary(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.Dashboard{
		Namespace: namespace,
		Quota:     quotaSnapshot,
		Contracts: summary,
	})
}

func (s *Server) handleCreateContract(w http.ResponseWriter, r *http.Request) {
	var input struct {
		NamespaceID string `json:"namespace_id"`
		Title       string `json:"title"`
		Signers     []struct {
			PartyName string `json:"party_name"`
			Email     string `json:"email"`
		} `json:"signers"`
		Content string `json:"content"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	exists, err := s.ns.Exists(input.NamespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	if !exists {
		writeError(w, http.StatusBadRequest, "namespace does not exist")
		return
	}
	signers := make([]doc.SignerInput, 0, len(input.Signers))
	for _, item := range input.Signers {
		signers = append(signers, doc.SignerInput{PartyName: item.PartyName, Email: item.Email})
	}
	contract, err := s.doc.Create(doc.CreateInput{
		NamespaceID: input.NamespaceID,
		Title:       input.Title,
		Signers:     signers,
		Content:     input.Content,
		Now:         now(),
	})
	if err != nil {
		handleError(w, err)
		return
	}
	_ = s.change.Append(contract.ID, "created", contract.Revision, 0, "contract created", now())
	_, _ = s.audit.Record(contract.ID, "contract.created", "success", input.Title, now())
	writeJSON(w, http.StatusCreated, contract)
}

func (s *Server) handleListContracts(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	keyword := r.URL.Query().Get("q")
	contracts, err := s.doc.Search(namespaceID, keyword)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, contracts)
}

func (s *Server) handleGetContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	contract, err := s.doc.Current(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, contract.Clone())
}

func (s *Server) handleContractContent(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	content, err := s.doc.Content(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, content)
}

func (s *Server) handleContractSigners(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	signers, err := s.sign.Signers(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, signers)
}

func (s *Server) handleSigningSummary(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	summary, err := s.sign.Summary(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) handleContractOrder(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	positions, err := s.order.Sequence(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, positions)
}

func (s *Server) handleSignerPosition(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	signerID := chi.URLParam(r, "signer_id")
	position, err := s.order.Position(contractID, signerID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, position)
}

func (s *Server) handleContractRevision(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	state, err := s.doc.RevisionState(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleContractRecords(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	records, err := s.sign.Records(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleContractLastRecord(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	record, err := s.sign.LastRecord(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleContractAcks(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	acks, err := s.doc.Acks(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, acks)
}

func (s *Server) handleContractCursor(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	cursor, err := s.sign.Cursor(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cursor)
}

func (s *Server) handleContractChanges(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	view, err := s.change.View(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleLatestChange(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	entry, err := s.change.Latest(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleContractAudit(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	events, err := s.audit.List(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleContractAuditRecent(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
			limit = 10
		}
	}
	events, err := s.audit.Recent(contractID, limit)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleLatestAudit(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	event, err := s.audit.Latest(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, event)
}

func (s *Server) handleContractTimeline(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	view, err := s.change.View(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	events, err := s.audit.List(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.TimelineFrom(view.Entries, events))
}

func (s *Server) handleArchiveDetail(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	record, err := s.archive.Latest(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleArchiveVerify(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	verification, err := s.archive.Verify(contractID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, verification)
}

func (s *Server) handleReviseContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	var input struct {
		Content string `json:"content"`
		Note    string `json:"note"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	contract, err := s.doc.Revise(contractID, input.Content, now())
	if err != nil {
		handleError(w, err)
		return
	}
	_ = s.change.Append(contract.ID, "revised", contract.Revision, contract.Revision-1, input.Note, now())
	_, _ = s.audit.Record(contract.ID, "contract.revised", "success", input.Note, now())
	writeJSON(w, http.StatusOK, contract)
}

func (s *Server) handleSignContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	var input struct {
		SignerID        string `json:"signer_id"`
		CertFingerprint string `json:"cert_fingerprint"`
		SignTime        string `json:"sign_time"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.order.Open(contractID); err != nil {
		handleError(w, err)
		return
	}
	record, err := s.sign.Do(contractID, input.SignerID, input.CertFingerprint, input.SignTime, now())
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleNotifySigners(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	var input struct {
		SignerID string `json:"signer_id"`
		InviteID string `json:"invite_id"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := s.sign.Notify(contractID, input.SignerID, input.InviteID, now())
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cursor)
}

func (s *Server) handleArchiveContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	var input struct {
		BatchNo string `json:"batch_no"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	record, err := s.archive.Run(contractID, input.BatchNo, now())
	if err != nil {
		handleError(w, err)
		return
	}
	_, _ = s.audit.Record(contractID, "contract.archived", "success", input.BatchNo, now())
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) handleArchiveRetry(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	var input struct {
		BatchNo string `json:"batch_no"`
	}
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	record, archived, err := s.archive.Retry(contractID, input.BatchNo, now())
	if err != nil {
		handleError(w, err)
		return
	}
	if !archived {
		writeJSON(w, http.StatusOK, map[string]any{"skipped": true, "record": nil})
		return
	}
	_, _ = s.audit.Record(contractID, "contract.archive-retried", "success", input.BatchNo, now())
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) handleArchiveRecords(w http.ResponseWriter, r *http.Request) {
	records, err := s.archive.List()
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleArchivePending(w http.ResponseWriter, r *http.Request) {
	namespaceID := r.URL.Query().Get("namespace_id")
	pending, err := s.archive.Pending(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pending)
}
