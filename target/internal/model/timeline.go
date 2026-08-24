package model

// TimelineEntry merges change journal rows and audit events into one
// chronological feed used by the audit console page.
type TimelineEntry struct {
	At     string `json:"at"`
	Kind   string `json:"kind"`
	Action string `json:"action"`
	Result string `json:"result,omitempty"`
	Note   string `json:"note,omitempty"`
	ID     string `json:"id"`
}

// TimelineFrom assembles a merged timeline; sources must already be sorted
// ascending by time.
func TimelineFrom(changes []ChangeEntry, events []AuditEvent) []TimelineEntry {
	entries := make([]TimelineEntry, 0, len(changes)+len(events))
	for _, change := range changes {
		entries = append(entries, TimelineEntry{
			At:     change.ChangedAt,
			Kind:   "change",
			Action: change.Action,
			Note:   change.Note,
			ID:     change.ID,
		})
	}
	for _, event := range events {
		entries = append(entries, TimelineEntry{
			At:     event.At,
			Kind:   "audit",
			Action: event.Action,
			Result: event.Result,
			Note:   event.Detail,
			ID:     event.ID,
		})
	}
	sortTimeline(entries)
	return entries
}

func sortTimeline(entries []TimelineEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].At < entries[j-1].At; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
