package verifycase

import (
	"os"
	"testing"

	"signflow/internal/audit"
	"signflow/internal/doc"
	"signflow/internal/order"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/sign"
	"signflow/internal/storage"
)

func TestSignCursorAfterAckDurable(t *testing.T) {
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
	if _, err := signSvc.Notify(contract.ID, contract.Signers[0].ID, "INV-001", now); err != nil {
		t.Fatal(err)
	}
	ackPath, err := fs.Path(cfg.AckFile(contract.ID)...)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(ackPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := signSvc.Notify(contract.ID, contract.Signers[1].ID, "INV-002", now); err == nil {
		t.Fatal("notify must fail when the acknowledgement cannot be written")
	}
	cursor, err := signSvc.Cursor(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cursor.LastAckID != "INV-001" {
		t.Fatalf("notification cursor advanced before the acknowledgement was durable: %s", cursor.LastAckID)
	}
}
