package aws

import (
	"io"
)

type S3Storage struct {
}

func (s *S3Storage) Create(docID string, rd io.Reader) error {
	panic("implement me")
}

func (s *S3Storage) Read(docID string) (io.Reader, error) {
	panic("implement me")
}

func (s *S3Storage) Update(docID string, rd io.Reader) error {
	panic("implement me")
}

func (s *S3Storage) Delete(docID string) error {
	panic("implement me")
}
