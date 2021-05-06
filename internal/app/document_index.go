package app

type DocumentIndex interface {
	Create(DocumentHeader) error
	GetByID(string) (DocumentHeader, error)
	Delete(string) error
}
