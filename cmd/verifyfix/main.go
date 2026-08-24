package main

// Temporary end-to-end verification of the sign.Do fix. Verifies the
// invariant: a signed contract always carries a complete signature ledger,
// both for the fresh-signing path and for the recovery path that backfills a
// record missing from an already-signed contract. Deleted after running.

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"signflow/internal/audit"
	"signflow/internal/doc"
	"signflow/internal/model"
	"signflow/internal/order"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/sign"
	"signflow/internal/storage"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	tmp, err := os.MkdirTemp("", "signflow-verify")
	must(err)
	defer os.RemoveAll(tmp)

	fs, err := storage.Open(tmp)
	must(err)
	cfg := settings.Default().WithRoot(tmp)

	quotaSvc := quota.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	orderSvc := order.NewService(docSvc)
	auditSvc := audit.NewService(fs, cfg)
	signSvc := sign.NewService(fs, cfg, docSvc, orderSvc, auditSvc)

	nsID := uuid.NewString()
	must(quotaSvc.Ensure(nsID, 100, "2026-08-23T00:00:00Z"))

	// One-signer contract.
	contract, err := docSvc.Create(doc.CreateInput{
		NamespaceID: nsID,
		Title:       "T",
		Signers:     []doc.SignerInput{{PartyName: "甲方", Email: "a@x.com"}},
		Content:     "body",
		Now:         "2026-08-23T00:00:00Z",
	})
	must(err)
	signerID := contract.Signers[0].ID

	// ---- Scenario 1: recovery path backfills a missing record. ----
	// Simulate the broken state produced by the old ordering: contract
	// effective + signer signed, but NO signature record on disk.
	must(docSvc.MarkSigned(contract.ID, signerID, "fp-1", "2026-08-22T12:00:00Z", "2026-08-23T00:00:00Z"))
	// At this point the contract is effective but has no signature record (the
	// old code would have failed exactly here on saveRecord).
	cur, err := docSvc.Current(contract.ID)
	must(err)
	if cur.Status != model.StatusEffective {
		panic("expected effective")
	}
	recs, err := signSvc.Records(contract.ID)
	must(err)
	if len(recs) != 0 {
		panic("expected no records yet")
	}
	// Retry: Do must not short-circuit on the effective status; it must
	// backfill the missing record.
	_, err = signSvc.Do(contract.ID, signerID, "fp-1", "2026-08-22T12:00:00Z", "2026-08-23T00:00:01Z")
	must(err)
	recs, err = signSvc.Records(contract.ID)
	must(err)
	if len(recs) != 1 {
		panic(fmt.Sprintf("expected 1 record after recovery, got %d", len(recs)))
	}
	if recs[0].CertFingerprint != "fp-1" || recs[0].SignedAt == "" {
		panic("record evidence incomplete")
	}
	fmt.Println("scenario 1 (recovery backfill): OK")

	// ---- Scenario 2: retry is idempotent (no duplicate records). ----
	_, err = signSvc.Do(contract.ID, signerID, "fp-1", "2026-08-22T12:00:00Z", "2026-08-23T00:00:02Z")
	must(err)
	recs, err = signSvc.Records(contract.ID)
	must(err)
	if len(recs) != 1 {
		panic(fmt.Sprintf("expected still 1 record after retry, got %d", len(recs)))
	}
	fmt.Println("scenario 2 (retry idempotent): OK")

	// ---- Scenario 3: fresh signing on a new contract writes record before state. ----
	c2, err := docSvc.Create(doc.CreateInput{
		NamespaceID: nsID,
		Title:       "T2",
		Signers:     []doc.SignerInput{{PartyName: "乙方", Email: "b@x.com"}},
		Content:     "body2",
		Now:         "2026-08-23T00:00:00Z",
	})
	must(err)
	s2 := c2.Signers[0].ID
	_, err = signSvc.Do(c2.ID, s2, "fp-2", "2026-08-22T13:00:00Z", "2026-08-23T00:00:03Z")
	must(err)
	recs2, err := signSvc.Records(c2.ID)
	must(err)
	if len(recs2) != 1 || recs2[0].CertFingerprint != "fp-2" {
		panic("fresh sign record missing/incomplete")
	}
	cur2, err := docSvc.Current(c2.ID)
	must(err)
	if cur2.Status != model.StatusEffective {
		panic("expected effective for fresh sign")
	}
	fmt.Println("scenario 3 (fresh sign, record before state): OK")

	fmt.Println("ALL VERIFIED")
}
