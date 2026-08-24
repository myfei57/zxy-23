package archive

import (
	"sort"
	"strings"

	"signflow/internal/model"
)

// List returns every archive ledger row across contracts ordered by time.
func (s *Service) List() ([]model.ArchiveRecord, error) {
	names, err := s.fs.List("archive", "records")
	if err != nil {
		return nil, err
	}
	var all []model.ArchiveRecord
	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		contractID := strings.TrimSuffix(name, ".json")
		records, err := s.Records(contractID)
		if err != nil {
			continue
		}
		all = append(all, records...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].ArchivedAt < all[j].ArchivedAt
	})
	return all, nil
}
