package model

// DocumentContent is the durable body of one contract revision.
type DocumentContent struct {
	ContractID string `json:"contract_id"`
	Revision   int    `json:"revision"`
	Content    string `json:"content"`
	UpdatedAt  string `json:"updated_at"`
}

// ChangeEntry is one immutable row in a contract's change journal.
type ChangeEntry struct {
	ID           string `json:"id"`
	ContractID   string `json:"contract_id"`
	Action       string `json:"action"`
	Revision     int    `json:"revision"`
	PrevRevision int    `json:"prev_revision"`
	Note         string `json:"note"`
	ChangedAt    string `json:"changed_at"`
}

// ChangeView is the page model for the change-history console page.
type ChangeView struct {
	ContractID      string        `json:"contract_id"`
	CurrentRevision int           `json:"current_revision"`
	Title           string        `json:"title"`
	Entries         []ChangeEntry `json:"entries"`
}
