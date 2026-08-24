package sign

import (
	"fmt"

	"signflow/internal/model"
)

// LastRecord returns the most recent durable signature record of a contract.
func (s *Service) LastRecord(contractID string) (*model.SignatureRecord, error) {
	records, err := s.Records(contractID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("sign: contract %s has no signature records", contractID)
	}
	return &records[len(records)-1], nil
}
