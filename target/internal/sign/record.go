package sign

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

// saveRecord durably appends one signature record to the contract ledger.
func (s *Service) saveRecord(record *model.SignatureRecord) error {
	return s.fs.AppendJSON(s.cfg.SignatureRecordFile(record.ContractID), record)
}

// Records returns every durable signature record of a contract in append order.
func (s *Service) Records(contractID string) ([]model.SignatureRecord, error) {
	var records []model.SignatureRecord
	if err := s.fs.ReadJSON(s.cfg.SignatureRecordFile(contractID), &records); err != nil {
		if storage.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return records, nil
}
