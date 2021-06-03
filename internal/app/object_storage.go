package app

import "io"

type ObjectStorage interface {
	Create(DocID, io.Reader) error
	Read(DocID) (io.ReadCloser, error)
	Update(DocID, io.Reader) error
	Delete(DocID) error
}
