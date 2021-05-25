package app

type DocumentHeader struct {
	ID   string
	Name string
	Size int64
}

type ACL struct {
	Permissions map[string]Permission
}

type Permission struct {
	Username string
	Read     bool
	Write    bool
	Delete   bool
	Share    bool
}
