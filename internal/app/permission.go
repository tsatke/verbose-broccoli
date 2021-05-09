package app

func (a *App) hasPerm(userID, docID string, m func(Permission) bool) bool {
	perm, err := a.permissions.Permissions(userID, docID)
	if err != nil {
		return false
	}
	return m(perm)
}

func (a *App) canRead(userID, docID string) bool {
	return a.hasPerm(userID, docID, func(p Permission) bool {
		return p.Read
	})
}

func (a *App) canWrite(userID, docID string) bool {
	return a.hasPerm(userID, docID, func(p Permission) bool {
		return p.Write
	})
}

func (a *App) canDelete(userID, docID string) bool {
	return a.hasPerm(userID, docID, func(p Permission) bool {
		return p.Delete
	})
}

func (a *App) canShare(userID, docID string) bool {
	return a.hasPerm(userID, docID, func(p Permission) bool {
		return p.Share
	})
}
