package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"signflow/internal/model"
)

// buildArchive serializes the contract snapshot and its document body into the
// durable archive artifact and returns the bytes, content hash and size.
func (s *Service) buildArchive(contract *model.Contract, content *model.DocumentContent) ([]byte, string, int64) {
	payload := struct {
		Contract *model.Contract        `json:"contract"`
		Content  *model.DocumentContent `json:"content"`
	}{
		Contract: contract,
		Content:  content,
	}
	encoded := mustMarshal(payload)
	sum := sha256.Sum256(encoded)
	return encoded, hex.EncodeToString(sum[:]), int64(len(encoded))
}

// writeArchiveFile durably writes the archive artifact of a contract.
func (s *Service) writeArchiveFile(contractID string, payload []byte) error {
	if err := s.fs.WriteFile(s.cfg.ArchiveFile(contractID), payload); err != nil {
		return fmt.Errorf("archive: persist artifact %s: %w", contractID, err)
	}
	return nil
}

func mustMarshal(value any) []byte {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	return payload
}
