package app

type DocumentIndex interface {
	Create(DocumentHeader, ACL) error
	GetByID(string) (DocumentHeader, error)
	Delete(string) error
	ACL(string) (ACL, error)
}
