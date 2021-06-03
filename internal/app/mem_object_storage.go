package app

import (
	"bytes"
	"fmt"
	"io"
)

type MemObjectStorage struct {
	data map[DocID][]byte
}

func NewMemObjectStorage() *MemObjectStorage {
	return &MemObjectStorage{
		data: map[DocID][]byte{},
	}
}

func (s *MemObjectStorage) Create(id DocID, rd io.Reader) error {
	_, ok := s.data[id]
	if ok {
		return fmt.Errorf("already exists")
	}

	data, err := io.ReadAll(rd)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	s.data[id] = data
	return nil
}

func (s *MemObjectStorage) Read(id DocID) (io.ReadCloser, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("does not exist")
	}
	return readCloserWrapper{bytes.NewReader(data)}, nil
}

func (s *MemObjectStorage) Update(id DocID, rd io.Reader) error {
	_, ok := s.data[id]
	if !ok {
		return fmt.Errorf("does not exist")
	}

	data, err := io.ReadAll(rd)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	s.data[id] = data
	return nil
}

func (s *MemObjectStorage) Delete(id DocID) error {
	delete(s.data, id)
	return nil
}

type readCloserWrapper struct {
	rd io.Reader
}

func (r readCloserWrapper) Read(p []byte) (n int, err error) {
	return r.rd.Read(p)
}

func (r readCloserWrapper) Close() error {
	return nil
}
