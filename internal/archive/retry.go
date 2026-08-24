package archive

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/doc"
	"signflow/internal/model"
)

// Retry resumes an interrupted archive batch. A contract that already owns an
// archive ledger row must be skipped so retry never duplicates an archive.
func (s *Service) Retry(contractID string, batchNo string, now string) (*model.ArchiveRecord, bool, error) {
	recorded, err := s.hasRecord(contractID)
	if err != nil {
		return nil, false, err
	}
	if recorded {
		return nil, false, nil
	}
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, false, err
	}
	if contract.Status != model.StatusEffective && contract.Status != model.StatusArchived {
		return nil, false, fmt.Errorf("%w: %s", doc.ErrNotEffective, contract.Status)
	}
	content, err := s.doc.Content(contractID)
	if err != nil {
		return nil, false, err
	}
	if content.Revision != contract.Revision {
		return nil, false, fmt.Errorf("archive: content revision %d does not match contract revision %d",
			content.Revision, contract.Revision)
	}
	payload, hash, size := s.buildArchive(contract, content)
	if err := s.writeArchiveFile(contractID, payload); err != nil {
		return nil, false, err
	}
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
		return nil, false, err
	}
	if contract.Status != model.StatusArchived {
		if err := s.doc.MarkArchived(contractID, now); err != nil {
			return nil, false, err
		}
	}
	return record, true, nil
}
