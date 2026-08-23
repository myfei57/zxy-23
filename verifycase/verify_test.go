package verifycase

import (
	"testing"

	"signflow/internal/archive"
	"signflow/internal/audit"
	"signflow/internal/doc"
	"signflow/internal/order"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/sign"
	"signflow/internal/storage"
)

func TestArchiveRetrySkipsRecordedContract(t *testing.T) {
	root := t.TempDir()
	cfg := settings.Default().WithRoot(root)
	fs, err := storage.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	quotaSvc := quota.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	orderSvc := order.NewService(docSvc)
	auditSvc := audit.NewService(fs, cfg)
	signSvc := sign.NewService(fs, cfg, docSvc, orderSvc, auditSvc)
	archiveSvc := archive.NewService(fs, cfg, docSvc)
	now := "2026-08-22T10:00:00Z"
	if err := quotaSvc.Ensure("ns-1", 10, now); err != nil {
		t.Fatal(err)
	}
	contract, err := docSvc.Create(doc.CreateInput{
		NamespaceID: "ns-1",
		Title:       "frame",
		Signers: []doc.SignerInput{
			{PartyName: "party-a", Email: "a@x.com"},
			{PartyName: "party-b", Email: "b@x.com"},
		},
		Content: "v1",
		Now:     now,
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := range contract.Signers {
		if _, err := signSvc.Do(contract.ID, contract.Signers[i].ID, "cert", now, now); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := archiveSvc.Run(contract.ID, "BATCH-01", now); err != nil {
		t.Fatal(err)
	}
	if _, archived, err := archiveSvc.Retry(contract.ID, "BATCH-02", now); err != nil {
		t.Fatal(err)
	} else if archived {
		t.Fatal("retry must skip a contract that already owns an archive record")
	}
	records, err := archiveSvc.Records(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("retry duplicated the archive ledger: rows=%d", len(records))
	}
}
