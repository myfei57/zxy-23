package audit

import (
	"fmt"

	"signflow/internal/model"
)

// Latest returns the most recent audit event of an entity.
func (s *Service) Latest(entityID string) (*model.AuditEvent, error) {
	events, err := s.List(entityID)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("audit: entity %s has no events", entityID)
	}
	return &events[len(events)-1], nil
}

// Recent returns up to limit most recent events, newest first.
func (s *Service) Recent(entityID string, limit int) ([]model.AuditEvent, error) {
	events, err := s.List(entityID)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}
	result := make([]model.AuditEvent, 0, limit)
	for i := len(events) - 1; i >= 0 && len(result) < limit; i-- {
		result = append(result, events[i])
	}
	return result, nil
}
