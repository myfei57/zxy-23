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

func TestAuditAfterSignDurable(t *testing.T) {
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
			{PartyName: "party-c", Email: "c@x.com"},
		},
		Content: "v1",
		Now:     now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := signSvc.Do(contract.ID, contract.Signers[0].ID, "cert-1", now, now); err != nil {
		t.Fatal(err)
	}
	recordPath, err := fs.Path(cfg.SignatureRecordFile(contract.ID)...)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(recordPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := signSvc.Do(contract.ID, contract.Signers[1].ID, "cert-2", now, now); err == nil {
		t.Fatal("signing must fail when the signature record cannot be written")
	}
	events, err := auditSvc.List(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Detail == contract.Signers[1].ID && event.Result == "success" && event.Action == "contract.signed" {
			t.Fatal("audit recorded a success event before the signature was durable")
		}
	}
}
