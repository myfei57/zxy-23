package model

// AuditEvent is one immutable audit row written by the audit component.
type AuditEvent struct {
	ID       string `json:"id"`
	EntityID string `json:"entity_id"`
	Action   string `json:"action"`
	Result   string `json:"result"`
	Detail   string `json:"detail"`
	At       string `json:"at"`
}
