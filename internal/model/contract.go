package model

// ContractStatus is the lifecycle state of a contract aggregate.
type ContractStatus string

const (
	StatusDraft     ContractStatus = "draft"
	StatusSigning   ContractStatus = "signing"
	StatusEffective ContractStatus = "effective"
	StatusArchived  ContractStatus = "archived"
)

// Contract is the central aggregate: version counter, document revision
// marker, signer list and lifecycle status are persisted together.
type Contract struct {
	ID          string         `json:"id"`
	NamespaceID string         `json:"namespace_id"`
	Title       string         `json:"title"`
	Status      ContractStatus `json:"status"`
	Revision    int            `json:"revision"`
	RevisedAt   string         `json:"revised_at,omitempty"`
	SignedAt    string         `json:"signed_at,omitempty"`
	EffectiveAt string         `json:"effective_at,omitempty"`
	ArchivedAt  string         `json:"archived_at,omitempty"`
	Signers     []Signer       `json:"signers"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// Clone returns a deep copy used for snapshots; mutating the copy never leaks
// into the persisted aggregate.
func (c *Contract) Clone() *Contract {
	if c == nil {
		return nil
	}
	copy := *c
	copy.Signers = make([]Signer, len(c.Signers))
	for i := range c.Signers {
		copy.Signers[i] = c.Signers[i]
	}
	return &copy
}

// SignerIndex returns the zero-based position of a signer or -1.
func (c *Contract) SignerIndex(signerID string) int {
	for i := range c.Signers {
		if c.Signers[i].ID == signerID {
			return i
		}
	}
	return -1
}

// PreviousSigner returns the signer that must finish before signerID may act.
func (c *Contract) PreviousSigner(signerID string) *Signer {
	index := c.SignerIndex(signerID)
	if index <= 0 {
		return nil
	}
	prev := c.Signers[index-1]
	return &prev
}

// NextSigner returns the first pending signer in signing order.
func (c *Contract) NextSigner() *Signer {
	for i := range c.Signers {
		if c.Signers[i].State != SignerSigned {
			return &c.Signers[i]
		}
	}
	return nil
}

// AllSigned reports whether every signer has a signed state.
func (c *Contract) AllSigned() bool {
	for i := range c.Signers {
		if c.Signers[i].State != SignerSigned {
			return false
		}
	}
	return len(c.Signers) > 0
}

// Signer returns a signer by ID.
func (c *Contract) Signer(signerID string) (*Signer, bool) {
	for i := range c.Signers {
		if c.Signers[i].ID == signerID {
			return &c.Signers[i], true
		}
	}
	return nil, false
}
