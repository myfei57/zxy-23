package model

// SignatureRecord is the durable evidence of one completed signature.
type SignatureRecord struct {
	ID              string `json:"id"`
	ContractID      string `json:"contract_id"`
	SignerID        string `json:"signer_id"`
	PartyName       string `json:"party_name"`
	CertFingerprint string `json:"cert_fingerprint"`
	SignedAt        string `json:"signed_at"`
	Status          string `json:"status"`
}

// AckRecord is the durable acknowledgement of one signing invitation.
type AckRecord struct {
	ID         string `json:"id"`
	ContractID string `json:"contract_id"`
	SignerID   string `json:"signer_id"`
	InviteID   string `json:"invite_id"`
	Delivered  bool   `json:"delivered"`
	AckedAt    string `json:"acked_at"`
}

// NotificationCursor records how far the invitation fan-out has progressed.
type NotificationCursor struct {
	ContractID   string `json:"contract_id"`
	LastAckID    string `json:"last_ack_id"`
	LastSignerID string `json:"last_signer_id"`
	UpdatedAt    string `json:"updated_at"`
}
