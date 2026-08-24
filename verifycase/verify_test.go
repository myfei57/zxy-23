package verifycase

import (
	"os"
	"testing"

	"signflow/internal/archive"
	"signflow/internal/audit"
	"signflow/internal/doc"
	"signflow/internal/model"
	"signflow/internal/order"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/sign"
	"signflow/internal/storage"
)

func TestArchiveFlagAfterFileDurable(t *testing.T) {
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
	archivePath, err := fs.Path(cfg.ArchiveFile(contract.ID)...)
	if err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile(cfg.ArchiveFile(contract.ID), []byte("placeholder")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(archivePath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := archiveSvc.Run(contract.ID, "BATCH-01", now); err == nil {
		t.Fatal("archive must fail when the archive artifact cannot be written")
	}
	current, err := docSvc.Current(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.Status == model.StatusArchived {
		t.Fatal("contract marked archived before the archive file was durable")
	}
}
