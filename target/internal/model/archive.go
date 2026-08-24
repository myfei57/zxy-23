package model

// ArchiveRecord is one row of the durable archive ledger.
type ArchiveRecord struct {
	ID          string `json:"id"`
	ContractID  string `json:"contract_id"`
	BatchNo     string `json:"batch_no"`
	ArchivePath string `json:"archive_path"`
	ContentHash string `json:"content_hash"`
	FileSize    int64  `json:"file_size"`
	ArchivedAt  string `json:"archived_at"`
}
