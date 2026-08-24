package doc

import (
	"signflow/internal/model"
	"signflow/internal/storage"
)

// Ack durably appends one signing invitation acknowledgement.
func (s *Service) Ack(contractID string, ack model.AckRecord) error {
	return s.fs.AppendJSON(s.cfg.AckFile(contractID), ack)
}

// Acks returns every acknowledgement of one contract in append order.
func (s *Service) Acks(contractID string) ([]model.AckRecord, error) {
	var acks []model.AckRecord
	if err := s.fs.ReadJSON(s.cfg.AckFile(contractID), &acks); err != nil {
		if storage.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return acks, nil
}
