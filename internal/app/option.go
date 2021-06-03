package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Option func(*App)

func WithCORSOrigins(origins ...string) Option {
	return func(a *App) {
		a.corsOrigins = append(a.corsOrigins, origins...)
	}
}

func WithRouter(r *gin.Engine) Option {
	return func(a *App) {
		a.router = r
	}
}

func WithObjectStorage(s ObjectStorage) Option {
	return func(a *App) {
		a.objects = s
	}
}

func WithDocumentRepo(r DocumentRepo) Option {
	return func(a *App) {
		a.documents = r
	}
}

func WithAuthService(s AuthService) Option {
	return func(a *App) {
		a.auth = s
	}
}

func WithHTTPServer(s *http.Server) Option {
	return func(a *App) {
		a.srv = s
	}
}

func WithLogger(log zerolog.Logger) Option {
	return func(a *App) {
		a.log = log
	}
}
