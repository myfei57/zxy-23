package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"signflow/internal/model"
)

// Verify re-hashes the durable archive artifact and compares it with the last
// ledger row, so tampering is detectable on demand.
func (s *Service) Verify(contractID string) (model.ArchiveVerification, error) {
	records, err := s.Records(contractID)
	if err != nil {
		return model.ArchiveVerification{}, err
	}
	verification := model.ArchiveVerification{
		ContractID: contractID,
		LedgerRows: len(records),
	}
	if len(records) == 0 {
		verification.Reason = "no archive ledger row"
		return verification, nil
	}
	latest := records[len(records)-1]
	verification.RecordedHash = latest.ContentHash
	payload, err := s.fs.ReadFile(s.cfg.ArchiveFile(contractID))
	if err != nil {
		verification.Reason = "archive artifact missing"
		return verification, nil
	}
	verification.FileExists = true
	size, modified, err := s.fs.Describe(s.cfg.ArchiveFile(contractID)...)
	if err == nil {
		verification.FileSize = size
		verification.ModifiedAt = modified.Format("2006-01-02T15:04:05Z07:00")
	}
	sum := sha256.Sum256(payload)
	verification.ActualHash = hex.EncodeToString(sum[:])
	verification.Valid = verification.RecordedHash == verification.ActualHash
	if !verification.Valid {
		verification.Reason = "content hash mismatch"
	}
	return verification, nil
}

// Latest returns the most recent archive ledger row of a contract.
func (s *Service) Latest(contractID string) (*model.ArchiveRecord, error) {
	records, err := s.Records(contractID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("archive: contract %s has no archive record", contractID)
	}
	return &records[len(records)-1], nil
}
