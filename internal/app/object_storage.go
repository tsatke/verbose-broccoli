package app

import "io"

type ObjectStorage interface {
	Create(string, io.Reader) error
	Read(string) (io.ReadCloser, error)
	Update(string, io.Reader) error
	Delete(string) error
}
