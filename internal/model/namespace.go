package model

// Namespace is the legal/contract namespace that owns contracts and quotas.
type Namespace struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LegalEntity string `json:"legal_entity"`
	CreatedAt   string `json:"created_at"`
}
