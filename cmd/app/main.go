package main

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/tsatke/verbose-broccoli/internal/app"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	u := app.NewMemUserService()
	_ = u.CreateUser("foo", "bar")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	a := app.New(
		lis,
		app.NewMemUserService(),
		app.NewMemObjectStorage(),
		app.NewMemDocumentIndex(),
		app.NewMemPermissionService(),
	)
	if err := a.Run(); err != nil {
		panic(err)
	}
}
