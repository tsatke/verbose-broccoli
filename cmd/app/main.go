package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/tsatke/verbose-broccoli/internal/app"
	appcfg "github.com/tsatke/verbose-broccoli/internal/app/config"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	c, err := appcfg.Load()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(
		c.GetString(appcfg.ListenerHost),
		c.GetString(appcfg.ListenerPort),
	))
	if err != nil {
		panic(err)
	}

	fmt.Println("listening on", lis.Addr().String())

	i, err := app.NewAuroraIndex(c)
	if err != nil {
		panic(err)
	}

	a := app.New(
		lis,
		app.NewS3Storage(c),
		i,
		app.NewCognitoService(c),
	)
	if err := a.Run(); err != nil {
		panic(err)
	}
}
