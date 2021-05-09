package app

type UserService interface {
	CredentialsValid(string, string) (bool, error)
	CreateUser(string, string) error
	UserID(string) (string, error)
}
