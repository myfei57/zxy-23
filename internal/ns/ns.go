package ns

import (
	"fmt"

	"github.com/google/uuid"

	"signflow/internal/model"
	"signflow/internal/settings"
	"signflow/internal/storage"
)

// Service manages legal/contract namespaces and their registry file.
type Service struct {
	fs  *storage.FS
	cfg *settings.Settings
}

// NewService builds the namespace service over the shared file store.
func NewService(fs *storage.FS, cfg *settings.Settings) *Service {
	return &Service{fs: fs, cfg: cfg}
}

// Create registers a new namespace and returns the persisted aggregate.
func (s *Service) Create(name string, legalEntity string, now string) (*model.Namespace, error) {
	namespace := &model.Namespace{
		ID:          uuid.NewString(),
		Name:        name,
		LegalEntity: legalEntity,
		CreatedAt:   now,
	}
	if err := s.fs.AppendJSON(s.cfg.NamespaceFile(), namespace); err != nil {
		return nil, fmt.Errorf("ns: persist namespace: %w", err)
	}
	return namespace, nil
}

// Get returns one namespace by ID.
func (s *Service) Get(id string) (*model.Namespace, error) {
	var namespaces []model.Namespace
	if err := s.fs.ReadJSON(s.cfg.NamespaceFile(), &namespaces); err != nil {
		return nil, err
	}
	for i := range namespaces {
		if namespaces[i].ID == id {
			return &namespaces[i], nil
		}
	}
	return nil, storage.ErrNotExist
}

// List returns every registered namespace in creation order.
func (s *Service) List() ([]model.Namespace, error) {
	var namespaces []model.Namespace
	if err := s.fs.ReadJSON(s.cfg.NamespaceFile(), &namespaces); err != nil {
		if storage.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return namespaces, nil
}

// Exists reports whether a namespace id is registered.
func (s *Service) Exists(id string) (bool, error) {
	_, err := s.Get(id)
	if err == nil {
		return true, nil
	}
	if storage.IsNotFound(err) {
		return false, nil
	}
	return false, err
}
