package app

import (
	"fmt"
)

type MemDocumentIndex struct {
	data map[string]DocumentHeader
	acls map[string]ACL
}

func (m *MemDocumentIndex) ACL(id string) (ACL, error) {
	if acl, ok := m.acls[id]; ok {
		return acl, nil
	}

	return ACL{}, fmt.Errorf("does not exist")
}

func NewMemDocumentIndex() *MemDocumentIndex {
	return &MemDocumentIndex{
		data: map[string]DocumentHeader{},
		acls: map[string]ACL{},
	}
}

func (m *MemDocumentIndex) Create(h DocumentHeader, acl ACL) error {
	if _, ok := m.data[h.ID]; ok {
		return fmt.Errorf("already exists")
	}

	m.data[h.ID] = h
	m.acls[h.ID] = acl

	return nil
}

func (m *MemDocumentIndex) GetByID(id string) (DocumentHeader, error) {
	if h, ok := m.data[id]; ok {
		return h, nil
	}

	return DocumentHeader{}, fmt.Errorf("does not exist")
}

func (m *MemDocumentIndex) Delete(id string) error {
	if _, ok := m.data[id]; !ok {
		return fmt.Errorf("does not exists")
	}

	delete(m.data, id)
	delete(m.acls, id)
	return nil
}
