package app

const (
	UserIDKey      = "UserID"
	UserIDTokenKey = "IDToken"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
