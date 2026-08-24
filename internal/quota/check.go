package quota

import (
	"errors"
	"fmt"
)

// ErrQuotaExceeded is returned when reserving capacity would exceed the limit.
var ErrQuotaExceeded = errors.New("quota: namespace capacity exceeded")

// Check verifies that delta additional contracts fit inside the namespace
// limit without mutating the ledger. The document flow must call it before any
// durable contract file is written.
func (s *Service) Check(namespaceID string, delta int) error {
	ledger, err := s.load(namespaceID)
	if err != nil {
		return fmt.Errorf("quota: check %s: %w", namespaceID, err)
	}
	if ledger.Used+delta > ledger.Limit {
		return fmt.Errorf("%w: used=%d limit=%d delta=%d", ErrQuotaExceeded, ledger.Used, ledger.Limit, delta)
	}
	return nil
}
