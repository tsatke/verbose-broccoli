package app

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type App struct {
	log         zerolog.Logger
	srv         *http.Server
	corsOrigins []string
	listener    net.Listener
	router      *gin.Engine
	objects     ObjectStorage
	index       DocumentIndex
	auth        AuthService
}

func New(lis net.Listener, opts ...Option) *App {
	a := &App{
		listener: lis,
		log:      zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.router == nil {
		a.router = gin.New()
		a.router.Use(gin.Recovery())
		a.router.Use(func(c *gin.Context) {
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}

			start := time.Now()
			c.Next()
			took := time.Since(start)

			var evt *zerolog.Event

			status := c.Writer.Status()
			switch {
			case status >= 500:
				evt = a.log.Error()
			case len(c.Errors) > 0:
				evt = a.log.Error()
			default:
				evt = a.log.Info()
			}

			if len(c.Errors) > 0 {
				evt = evt.Strs("errors", c.Errors.Errors())
			}

			evt = evt.
				Str("method", c.Request.Method).
				Str("path", path).
				Int("status", c.Writer.Status()).
				Str("ip", c.ClientIP()).
				Stringer("took", took)

			evt.Msg("request")
		})
	}
	if a.objects == nil {
		a.objects = NewMemObjectStorage()
	}
	if a.index == nil {
		a.index = NewMemDocumentIndex()
	}
	if a.auth == nil {
		a.auth = NewMemAuthService()
	}
	if a.srv == nil {
		a.srv = &http.Server{
			Handler:           a.router,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
	}

	a.setupCORS()
	a.setupRoutes()
	return a
}

func (a *App) Run() error {
	a.log.
		Info().
		IPAddr("host", a.listener.Addr().(*net.TCPAddr).IP).
		Int("port", a.listener.Addr().(*net.TCPAddr).Port).
		Msg("run server")
	if err := a.srv.Serve(a.listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Close() error {
	return a.srv.Close()
}
