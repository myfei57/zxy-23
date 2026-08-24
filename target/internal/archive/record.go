package archive

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

// appendRecord durably appends one row to the archive ledger of a contract.
func (s *Service) appendRecord(record *model.ArchiveRecord) error {
	return s.fs.AppendJSON(s.cfg.ArchiveRecordFile(record.ContractID), record)
}

// Records returns every archive ledger row of a contract in append order.
func (s *Service) Records(contractID string) ([]model.ArchiveRecord, error) {
	var records []model.ArchiveRecord
	if err := s.fs.ReadJSON(s.cfg.ArchiveRecordFile(contractID), &records); err != nil {
		if storage.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return records, nil
}

// hasRecord reports whether a contract already owns an archive ledger row.
func (s *Service) hasRecord(contractID string) (bool, error) {
	records, err := s.Records(contractID)
	if err != nil {
		return false, err
	}
	return len(records) > 0, nil
}
