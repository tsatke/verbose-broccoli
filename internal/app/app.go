package app

import (
	"net"

	"github.com/gin-gonic/gin"
)

type App struct {
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
	return a
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
