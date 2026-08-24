package sign

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/doc"
	"signflow/internal/model"
)

// Do completes one signature: the signature record must be durable before the
// contract is marked signed, and the audit success event is written last.
//
// Ordering the evidence before the state guarantees the invariant that a
// signed contract always carries a complete signature ledger: if the record
// write fails, the contract is never advanced to signed, so a receipt-timeout
// retry can simply run again instead of being blocked by the already-signed
// state. Should a signed contract ever lack a signer's record anyway (for
// example pre-existing data written by an older version), the same retry path
// backfills the missing record instead of short-circuiting on ErrClosed.
func (s *Service) Do(contractID string, signerID string, cert string, signTime string, now string) (*model.SignatureRecord, error) {
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	signer, ok := contract.Signer(signerID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", doc.ErrSignerUnknown, signerID)
	}
	if cert == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingCertificate, signerID)
	}
	// A signer already marked signed may still be missing its durable record
	// if a previous write failed after the state advanced. In that case the
	// signing order no longer applies (the slot is taken) and we must complete
	// the ledger rather than reject the retry as closed.
	alreadySigned := signer.State == model.SignerSigned
	if !alreadySigned {
		if contract.Status == model.StatusArchived || contract.Status == model.StatusEffective {
			return nil, fmt.Errorf("%w: %s", ErrClosed, contract.Status)
		}
		if err := s.order.Check(contractID, signerID); err != nil {
			return nil, err
		}
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
	// Persist the evidence first, so a write failure leaves the contract
	// untouched and retryable. The append is idempotent per signer, which lets
	// a retry safely run even when an earlier attempt already wrote the record.
	if err := s.saveRecord(record); err != nil {
		return nil, err
	}
	if err := s.doc.MarkSigned(contractID, signerID, cert, signTime, now); err != nil {
		return nil, err
	}
	if _, err := s.audit.Record(contractID, "contract.signed", "success", signerID, now); err != nil {
		return nil, err
	}
	return record, nil
}
