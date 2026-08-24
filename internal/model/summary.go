package model

// ContractSummary counts contracts per lifecycle status inside a namespace.
type ContractSummary struct {
	NamespaceID string `json:"namespace_id"`
	Draft       int    `json:"draft"`
	Signing     int    `json:"signing"`
	Effective   int    `json:"effective"`
	Archived    int    `json:"archived"`
	Total       int    `json:"total"`
}

// SigningSummary is the signing-page read model of one contract.
type SigningSummary struct {
	ContractID string  `json:"contract_id"`
	Total      int     `json:"total"`
	Signed     int     `json:"signed"`
	Pending    int     `json:"pending"`
	NextSigner *Signer `json:"next_signer,omitempty"`
	Complete   bool    `json:"complete"`
}

// OrderPosition is one row of the signing-sequence read model.
type OrderPosition struct {
	SignerID  string      `json:"signer_id"`
	PartyName string      `json:"party_name"`
	Position  int         `json:"position"`
	Total     int         `json:"total"`
	State     SignerState `json:"state"`
	IsNext    bool        `json:"is_next"`
	BlockedBy string      `json:"blocked_by,omitempty"`
}

// ArchiveVerification is the tamper-check result of one archive artifact.
type ArchiveVerification struct {
	ContractID   string `json:"contract_id"`
	RecordedHash string `json:"recorded_hash,omitempty"`
	ActualHash   string `json:"actual_hash,omitempty"`
	FileExists   bool   `json:"file_exists"`
	FileSize     int64  `json:"file_size,omitempty"`
	ModifiedAt   string `json:"modified_at,omitempty"`
	LedgerRows   int    `json:"ledger_rows"`
	Valid        bool   `json:"valid"`
	Reason       string `json:"reason,omitempty"`
}

// Dashboard is the namespace overview returned to the contracts page.
type Dashboard struct {
	Namespace *Namespace      `json:"namespace"`
	Quota     QuotaSnapshot   `json:"quota"`
	Contracts ContractSummary `json:"contracts"`
}
