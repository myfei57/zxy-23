package audit

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/model"
)

// Record appends one audit event and returns it. The signing flow must call it
// only after the signature record is durable.
func (s *Service) Record(entityID string, action string, result string, detail string, at string) (*model.AuditEvent, error) {
	if len(detail) > 500 {
		detail = detail[:500]
	}
	event := &model.AuditEvent{
		ID:       uuid.NewString(),
		EntityID: entityID,
		Action:   action,
		Result:   result,
		Detail:   detail,
		At:       at,
	}
	if err := s.fs.AppendJSON(s.cfg.AuditEventFile(entityID), event); err != nil {
		return nil, fmt.Errorf("audit: persist event: %w", err)
	}
	return event, nil
}
