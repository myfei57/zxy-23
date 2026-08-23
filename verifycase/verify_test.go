package verifycase

import (
	"os"
	"testing"

	"signflow/internal/doc"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

func TestRevisionMarkAfterContentDurable(t *testing.T) {
	root := t.TempDir()
	cfg := settings.Default().WithRoot(root)
	fs, err := storage.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	quotaSvc := quota.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	now := "2026-08-22T10:00:00Z"
	if err := quotaSvc.Ensure("ns-1", 10, now); err != nil {
		t.Fatal(err)
	}
	contract, err := docSvc.Create(doc.CreateInput{
		NamespaceID: "ns-1",
		Title:       "frame",
		Signers: []doc.SignerInput{
			{PartyName: "party-a", Email: "a@x.com"},
		},
		Content: "v1",
		Now:     now,
	})
	if err != nil {
		t.Fatal(err)
	}
	contentPath, err := fs.Path(cfg.ContentFile(contract.ID)...)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(contentPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := docSvc.Revise(contract.ID, "v2", now); err == nil {
		t.Fatal("revision must fail when the new content cannot be written")
	}
	current, err := docSvc.Current(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.RevisedAt != "" {
		t.Fatalf("revision marker persisted before content was durable: revised_at=%s", current.RevisedAt)
	}
}
