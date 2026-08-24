package verifycase

import (
	"testing"

	"signflow/internal/change"
	"signflow/internal/doc"
	"signflow/internal/quota"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

func TestChangeViewUsesCurrentDoc(t *testing.T) {
	root := t.TempDir()
	cfg := settings.Default().WithRoot(root)
	fs, err := storage.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	quotaSvc := quota.NewService(fs, cfg)
	docSvc := doc.NewService(fs, cfg, quotaSvc)
	changeSvc := change.NewService(fs, cfg, docSvc)
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
	if err := changeSvc.Append(contract.ID, "created", 1, 0, "created", now); err != nil {
		t.Fatal(err)
	}
	first, err := changeSvc.View(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if first.CurrentRevision != 1 {
		t.Fatalf("unexpected first view revision: %d", first.CurrentRevision)
	}
	if _, err := docSvc.Revise(contract.ID, "v2", now); err != nil {
		t.Fatal(err)
	}
	if err := changeSvc.Append(contract.ID, "revised", 2, 1, "revised", now); err != nil {
		t.Fatal(err)
	}
	second, err := changeSvc.View(contract.ID)
	if err != nil {
		t.Fatal(err)
	}
	if second.CurrentRevision != 2 {
		t.Fatalf("change history kept showing the stale revision: %d", second.CurrentRevision)
	}
}
