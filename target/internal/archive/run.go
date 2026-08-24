package archive

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"signflow/internal/doc"
	"signflow/internal/model"
)

// Run archives one effective contract. The archive artifact must be durable
// before the ledger row and the archived marker are persisted.
func (s *Service) Run(contractID string, batchNo string, now string) (*model.ArchiveRecord, error) {
	if strings.TrimSpace(batchNo) == "" {
		return nil, fmt.Errorf("archive: batch number is required")
	}
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	if contract.Status != model.StatusEffective {
		return nil, fmt.Errorf("%w: %s", doc.ErrNotEffective, contract.Status)
	}
	content, err := s.doc.Content(contractID)
	if err != nil {
		return nil, err
	}
	if content.Revision != contract.Revision {
		return nil, fmt.Errorf("archive: content revision %d does not match contract revision %d",
			content.Revision, contract.Revision)
	}
	payload, hash, size := s.buildArchive(contract, content)
	record := &model.ArchiveRecord{
		ID:          uuid.NewString(),
		ContractID:  contractID,
		BatchNo:     batchNo,
		ArchivePath: "archive/files/" + contractID + ".json",
		ContentHash: hash,
		FileSize:    size,
		ArchivedAt:  now,
	}
	if err := s.appendRecord(record); err != nil {
		return nil, err
	}
	if err := s.doc.MarkArchived(contractID, now); err != nil {
		return nil, err
	}
	if err := s.writeArchiveFile(contractID, payload); err != nil {
		return nil, err
	}
	return record, nil
}
