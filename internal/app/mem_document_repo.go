package app

import (
	"fmt"
)

type MemDocumentRepo struct {
	data map[DocID]DocumentHeader
	acls map[DocID]ACL
}

func (m *MemDocumentRepo) ACL(id DocID) (ACL, error) {
	if acl, ok := m.acls[id]; ok {
		return acl, nil
	}

	return ACL{}, fmt.Errorf("does not exist")
}

func NewMemDocumentRepo() *MemDocumentRepo {
	return &MemDocumentRepo{
		data: map[DocID]DocumentHeader{},
		acls: map[DocID]ACL{},
	}
}

func (m *MemDocumentRepo) Create(h DocumentHeader, acl ACL) error {
	if _, ok := m.data[h.ID]; ok {
		return fmt.Errorf("already exists")
	}
	return m.Update(h, acl)
}

func (m *MemDocumentRepo) Update(h DocumentHeader, acl ACL) error {
	m.data[h.ID] = h
	m.acls[h.ID] = acl

	return nil
}

func (m *MemDocumentRepo) Get(id DocID) (DocumentHeader, error) {
	if h, ok := m.data[id]; ok {
		return h, nil
	}

	return DocumentHeader{}, fmt.Errorf("does not exist")
}

func (m *MemDocumentRepo) Delete(id DocID) error {
	if _, ok := m.data[id]; !ok {
		return fmt.Errorf("does not exists")
	}

	delete(m.data, id)
	delete(m.acls, id)
	return nil
}
