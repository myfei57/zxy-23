package verifycase

import (
	"testing"

	"signflow/internal/doc"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

func TestContractQuotaRejectsBeforeCreate(t *testing.T) {
	root := t.TempDir()
	cfg := settings.Default().WithRoot(root)
	fs, err := storage.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	quotaSvc := quota.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	now := "2026-08-22T10:00:00Z"
	if err := quotaSvc.Ensure("ns-1", 1, now); err != nil {
		t.Fatal(err)
	}
	input := doc.CreateInput{
		NamespaceID: "ns-1",
		Title:       "frame",
		Signers: []doc.SignerInput{
			{PartyName: "party-a", Email: "a@x.com"},
		},
		Content: "v1",
		Now:     now,
	}
	if _, err := docSvc.Create(input); err != nil {
		t.Fatal(err)
	}
	over := input
	over.Title = "over-quota"
	if _, err := docSvc.Create(over); err == nil {
		t.Fatal("over-quota create must fail")
	}
	names, err := fs.List("doc", "contracts")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 {
		t.Fatalf("over-quota contract consumed storage: contract files=%d", len(names))
	}
}
