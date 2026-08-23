package model

// SignerState tracks one party inside a contract's signing sequence.
type SignerState string

const (
	SignerPending SignerState = "pending"
	SignerSigned  SignerState = "signed"
)

// Signer is one party of a contract with an explicit order position.
type Signer struct {
	ID              string      `json:"id"`
	ContractID      string      `json:"contract_id"`
	PartyName       string      `json:"party_name"`
	Email           string      `json:"email"`
	Order           int         `json:"order"`
	State           SignerState `json:"state"`
	SignedAt        string      `json:"signed_at,omitempty"`
	CertFingerprint string      `json:"cert_fingerprint,omitempty"`
	SignTime        string      `json:"sign_time,omitempty"`
}

// MarkSigned mutates the signer into the signed state and stamps the evidence.
func (s *Signer) MarkSigned(cert string, signTime string, at string) {
	s.State = SignerSigned
	s.SignedAt = at
	s.CertFingerprint = cert
	s.SignTime = signTime
}
