package app

import "time"

type (
	DocID string

	DocumentHeader struct {
		ID      DocID
		Name    string
		Owner   string
		Created time.Time
		Updated time.Time
	}

	ACL struct {
		Permissions map[string]Permission
	}

	Permission struct {
		Username string
		Read     bool
		Write    bool
		Delete   bool
		Share    bool
	}
)

type DocumentRepo interface {
	Create(DocumentHeader, ACL) error
	Get(DocID) (DocumentHeader, error)
	Delete(DocID) error
	ACL(DocID) (ACL, error)
}
