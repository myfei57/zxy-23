package sign

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

// saveRecord durably persists one signature record on the contract ledger.
//
// It is idempotent per signer: at most one record is kept for each signer, so
// a receipt-timeout retry that re-enters Do never leaves duplicate evidence
// rows and a recovery retry can backfill a record that an earlier attempt
// never wrote. This is what makes the "signed implies complete ledger"
// invariant repairable instead of merely avoidable.
func (s *Service) saveRecord(record *model.SignatureRecord) error {
	records, err := s.Records(record.ContractID)
	if err != nil {
		return err
	}
	for i := range records {
		if records[i].SignerID == record.SignerID {
			records[i] = *record
			return s.fs.WriteJSON(s.cfg.SignatureRecordFile(record.ContractID), records)
		}
	}
	records = append(records, *record)
	return s.fs.WriteJSON(s.cfg.SignatureRecordFile(record.ContractID), records)
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
