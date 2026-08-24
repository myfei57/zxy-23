package audit

import (
	"strings"

	"signflow/internal/storage"
)

// Total counts every audit event across all entities.
func (s *Service) Total() (int, error) {
	names, err := s.fs.List("audit", "events")
	if err != nil {
		return 0, err
	}
	total := 0
	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		entityID := strings.TrimSuffix(name, ".json")
		events, err := s.List(entityID)
		if err != nil {
			if storage.IsNotFound(err) {
				continue
			}
			return 0, err
		}
		total += len(events)
	}
	return total, nil
}
