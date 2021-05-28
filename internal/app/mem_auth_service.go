package app

import "github.com/google/uuid"

type MemAuthService struct {
	data   map[string]string
	tokens map[string]string
}

func NewMemAuthService() *MemAuthService {
	return &MemAuthService{
		data:   map[string]string{},
		tokens: map[string]string{},
	}
}

func (m *MemAuthService) Login(user, pass string) (LoginResult, error) {
	if m.data[user] == pass {
		token := uuid.New().String()
		m.tokens[token] = user
		return LoginResult{
			Success: true,
			Token:   token,
		}, nil
	}
	return LoginResult{
		Success: false,
	}, nil
}

func (m *MemAuthService) AnswerChallenge(user, challenge, payload string) (LoginResult, error) {
	return LoginResult{}, nil
}

func (m *MemAuthService) TokenValid(s string) bool {
	_, ok := m.tokens[s]
	return ok
}
