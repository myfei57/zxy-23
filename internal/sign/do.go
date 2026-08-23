package sign

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/doc"
	"signflow/internal/model"
)

// Do completes one signature: the signature record must be durable before the
// contract is marked signed, and the audit success event is written last.
func (s *Service) Do(contractID string, signerID string, cert string, signTime string, now string) (*model.SignatureRecord, error) {
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	if contract.Status == model.StatusArchived || contract.Status == model.StatusEffective {
		return nil, fmt.Errorf("%w: %s", ErrClosed, contract.Status)
	}
	signer, ok := contract.Signer(signerID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", doc.ErrSignerUnknown, signerID)
	}
	if cert == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCertificate, signerID)
	}
	if err := s.order.Check(contractID, signerID); err != nil {
		return nil, err
	}
	record := &model.SignatureRecord{
		ID:              uuid.NewString(),
		ContractID:      contractID,
		SignerID:        signerID,
		PartyName:       signer.PartyName,
		CertFingerprint: cert,
		SignedAt:        now,
		Status:          "completed",
	}
	if _, err := s.audit.Record(contractID, "contract.signed", "success", signerID, now); err != nil {
		return nil, err
	}
	if err := s.saveRecord(record); err != nil {
		return nil, err
	}
	if err := s.doc.MarkSigned(contractID, signerID, cert, signTime, now); err != nil {
		return nil, err
	}
	return record, nil
}
