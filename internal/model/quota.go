package model

// QuotaLedger is the durable capacity ledger of one namespace.
type QuotaLedger struct {
	NamespaceID string `json:"namespace_id"`
	Limit       int    `json:"limit"`
	Used        int    `json:"used"`
	UpdatedAt   string `json:"updated_at"`
}

// QuotaSnapshot is the read model returned to the quota console API.
type QuotaSnapshot struct {
	NamespaceID string `json:"namespace_id"`
	Limit       int    `json:"limit"`
	Used        int    `json:"used"`
	Remaining   int    `json:"remaining"`
}

// Snapshot converts a ledger into the API read model.
func (l *QuotaLedger) Snapshot() QuotaSnapshot {
	remaining := l.Limit - l.Used
	if remaining < 0 {
		remaining = 0
	}
	return QuotaSnapshot{
		NamespaceID: l.NamespaceID,
		Limit:       l.Limit,
		Used:        l.Used,
		Remaining:   remaining,
	}
}
