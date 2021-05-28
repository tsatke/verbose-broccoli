package main

import (
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

	// i, err := app.NewAuroraIndex(c)
	// if err != nil {
	// 	panic(err)
	// }

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Logger()

	a := app.New(lis,
		app.WithLogger(log),
		app.WithObjectStorage(app.NewS3Storage(c)),
		// app.WithDocumentIndex(i),
		app.WithAuthService(app.NewCognitoService(c)),
	)
	if err := a.Run(); err != nil {
		panic(err)
	}
}
