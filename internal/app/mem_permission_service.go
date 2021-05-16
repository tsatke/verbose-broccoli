package app

import "fmt"

type MemPermissionService struct {
	data map[[2]string]Permission
}

func NewMemPermissionService() *MemPermissionService {
	return &MemPermissionService{
		data: map[[2]string]Permission{},
	}
}

func (s *MemPermissionService) Permissions(userID, docID string) (Permission, error) {
	p, ok := s.data[[2]string{userID, docID}]
	if !ok {
		return Permission{}, fmt.Errorf("does not exist")
	}
	return p, nil
}

func (s *MemPermissionService) Create(p Permission) error {
	s.data[[2]string{p.User, p.Document}] = p
	return nil
}

func (s *MemPermissionService) Delete(userID, docID string) error {
	delete(s.data, [2]string{userID, docID})
	return nil
}

func (s *MemPermissionService) All(userID string) ([]Permission, error) {
	var p []Permission

	for strings, permission := range s.data {
		if strings[0] == userID {
			p = append(p, permission)
		}
	}

	return p, nil
}
