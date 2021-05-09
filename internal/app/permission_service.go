package app

type Permission struct {
	User     string
	Document string

	Read   bool
	Write  bool
	Delete bool
	Share  bool
}

type PermissionService interface {
	Permissions(userID, docID string) (Permission, error)
	Create(Permission) error
	Delete(userID, docID string) error
	All(userID string) ([]Permission, error)
}
