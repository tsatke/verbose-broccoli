package app

type MemUserService struct {
	userIDs   map[string]string
	passwords map[string]string
}

func NewMemUserService() *MemUserService {
	return &MemUserService{
		userIDs:   map[string]string{},
		passwords: map[string]string{},
	}
}

func (m *MemUserService) CredentialsValid(username, password string) (bool, error) {
	return m.passwords[username] == password, nil
}

func (m *MemUserService) CreateUser(username, password string) error {
	m.userIDs[username] = username
	m.passwords[username] = password
	return nil
}

func (m *MemUserService) UserID(username string) (string, error) {
	return m.userIDs[username], nil
}
