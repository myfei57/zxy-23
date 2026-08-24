package settings

// Settings carries the file-persistence layout used by every store in the
// platform. All relative paths are anchored below Root so a service can run
// against a scratch directory or a production data directory without changes.
type Settings struct {
	Root string
}

// Default returns the settings used by the command entry point before the
// -root flag overrides the data directory.
func Default() *Settings {
	return &Settings{Root: "signflow-data"}
}

// WithRoot returns a copy of the settings with a different persistence root.
func (s *Settings) WithRoot(root string) *Settings {
	clone := *s
	clone.Root = root
	return &clone
}

// NamespaceFile returns the relative path of the namespace registry.
func (s *Settings) NamespaceFile() []string {
	return []string{"ns", "namespaces.json"}
}

// ContractFile returns the relative path of one contract aggregate.
func (s *Settings) ContractFile(contractID string) []string {
	return []string{"doc", "contracts", contractID + ".json"}
}

// ContentFile returns the relative path of the latest durable document content.
func (s *Settings) ContentFile(contractID string) []string {
	return []string{"doc", "content", contractID + ".json"}
}

// AckFile returns the relative path of one contract's invite acknowledgements.
func (s *Settings) AckFile(contractID string) []string {
	return []string{"doc", "ack", contractID + ".json"}
}

// SignatureRecordFile returns the relative path of one contract's signature
// records; the array is appended to after each durable signature.
func (s *Settings) SignatureRecordFile(contractID string) []string {
	return []string{"sign", "records", contractID + ".json"}
}

// NotificationCursorFile returns the relative path of one contract's signing
// notification cursor.
func (s *Settings) NotificationCursorFile(contractID string) []string {
	return []string{"sign", "cursor", contractID + ".json"}
}

// ArchiveFile returns the relative path of the durable archive artifact.
func (s *Settings) ArchiveFile(contractID string) []string {
	return []string{"archive", "files", contractID + ".json"}
}

// ArchiveRecordFile returns the relative path of one contract's archive ledger.
func (s *Settings) ArchiveRecordFile(contractID string) []string {
	return []string{"archive", "records", contractID + ".json"}
}

// ChangeJournalFile returns the relative path of one contract's change journal.
func (s *Settings) ChangeJournalFile(contractID string) []string {
	return []string{"change", "journal", contractID + ".json"}
}

// AuditEventFile returns the relative path of one entity's audit events.
func (s *Settings) AuditEventFile(entityID string) []string {
	return []string{"audit", "events", entityID + ".json"}
}

// QuotaLedgerFile returns the relative path of one namespace's quota ledger.
func (s *Settings) QuotaLedgerFile(namespaceID string) []string {
	return []string{"quota", "ledger", namespaceID + ".json"}
}
