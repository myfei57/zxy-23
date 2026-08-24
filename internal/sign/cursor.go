package sign

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"signflow/internal/model"
	"signflow/internal/storage"
)

// Notify delivers one signing invitation. The acknowledgement must be durable
// before the notification cursor advances, so a failed acknowledgement never
// skips a signer on restart.
func (s *Service) Notify(contractID string, signerID string, inviteID string, now string) (*model.NotificationCursor, error) {
	contract, err := s.doc.Current(contractID)
	if err != nil {
		return nil, err
	}
	if _, ok := contract.Signer(signerID); !ok {
		return nil, fmt.Errorf("sign: signer %s is not part of contract %s", signerID, contractID)
	}
	if strings.TrimSpace(inviteID) == "" {
		return nil, fmt.Errorf("sign: invite id is required")
	}
	ack := &model.AckRecord{
		ID:         uuid.NewString(),
		ContractID: contractID,
		SignerID:   signerID,
		InviteID:   inviteID,
		Delivered:  true,
		AckedAt:    now,
	}
	if err := s.doc.Ack(contractID, *ack); err != nil {
		return nil, fmt.Errorf("sign: persist acknowledgement: %w", err)
	}
	return s.saveCursor(contractID, inviteID, signerID, now)
}

// saveCursor durably advances the notification cursor.
func (s *Service) saveCursor(contractID string, lastAckID string, lastSignerID string, now string) (*model.NotificationCursor, error) {
	cursor := &model.NotificationCursor{
		ContractID:   contractID,
		LastAckID:    lastAckID,
		LastSignerID: lastSignerID,
		UpdatedAt:    now,
	}
	if err := s.fs.WriteJSON(s.cfg.NotificationCursorFile(contractID), cursor); err != nil {
		return nil, fmt.Errorf("sign: persist notification cursor: %w", err)
	}
	return cursor, nil
}

// Cursor returns the last durable notification cursor of a contract.
func (s *Service) Cursor(contractID string) (*model.NotificationCursor, error) {
	var cursor model.NotificationCursor
	if err := s.fs.ReadJSON(s.cfg.NotificationCursorFile(contractID), &cursor); err != nil {
		if storage.IsNotFound(err) {
			return &model.NotificationCursor{ContractID: contractID}, nil
		}
		return nil, err
	}
	return &cursor, nil
}
