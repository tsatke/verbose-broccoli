package mem

import (
	"fmt"

	"github.com/tsatke/verbose-broccoli/internal/app"
)

type DocumentIndex struct {
	data map[string]app.DocumentHeader
}

func (m *DocumentIndex) Create(h app.DocumentHeader) error {
	if _, ok := m.data[h.ID]; ok {
		return fmt.Errorf("already exists")
	}

	m.data[h.ID] = h

	return nil
}

func (m *DocumentIndex) GetByID(id string) (app.DocumentHeader, error) {
	if h, ok := m.data[id]; ok {
		return h, nil
	}

	return app.DocumentHeader{}, fmt.Errorf("does not exist")
}

func (m *DocumentIndex) Delete(id string) error {
	if _, ok := m.data[id]; !ok {
		return fmt.Errorf("does not exists")
	}

	delete(m.data, id)
	return nil
}
