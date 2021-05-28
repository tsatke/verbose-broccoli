package app

//go:generate mockery --inpackage --testonly --case snake --name AuthService --filename auth_service_mock_test.go

type LoginResult struct {
	Success   bool
	Challenge string
	Token     string
}

type AuthService interface {
	Login(user, pass string) (LoginResult, error)
	AnswerChallenge(user, challenge, payload string) (LoginResult, error)
	TokenValid(string) bool
}
