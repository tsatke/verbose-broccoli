package mem

type UserService struct {
	userIDs   map[string]string
	passwords map[string]string
}

func NewUserService() *UserService {
	return &UserService{
		userIDs:   map[string]string{},
		passwords: map[string]string{},
	}
}

func (m *UserService) CredentialsValid(username string, password string) (bool, error) {
	return m.passwords[username] == password, nil
}

func (m *UserService) CreateUser(username string, password string) error {
	m.userIDs[username] = username
	m.passwords[username] = password
	return nil
}

func (m *UserService) UserID(username string) (string, error) {
	return m.userIDs[username], nil
}
