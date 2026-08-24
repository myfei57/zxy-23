package ns

// Count returns the number of registered namespaces.
func (s *Service) Count() (int, error) {
	namespaces, err := s.List()
	if err != nil {
		return 0, err
	}
	return len(namespaces), nil
}
