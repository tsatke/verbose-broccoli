package app

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	srv         *http.Server
	listener    net.Listener
	router      *gin.Engine
	objects     ObjectStorage
	index       DocumentIndex
	permissions PermissionService
	users       UserService
}

func New(l net.Listener, u UserService, o ObjectStorage, i DocumentIndex, p PermissionService) *App {
	a := &App{
		listener:    l,
		router:      gin.Default(),
		objects:     o,
		index:       i,
		permissions: p,
		users:       u,
	}
	a.setupRoutes()
	a.srv = &http.Server{
		Handler:           a.router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return a
}

func (a *App) Run() error {
	if err := a.srv.Serve(a.listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Close() error {
	return a.srv.Close()
}
