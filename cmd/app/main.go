package main

import (
	"fmt"
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
		fatal(err)
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(
		c.GetString(appcfg.ListenerHost),
		c.GetString(appcfg.ListenerPort),
	))
	if err != nil {
		fatal(err)
	}

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Logger()

	p, err := app.NewPostgresDatabaseProvider(log, c, true)
	if err != nil {
		fatal(err)
	}

	a := app.New(lis,
		app.WithLogger(log),
		app.WithObjectStorage(app.NewS3Storage(c)),
		app.WithDocumentRepo(app.NewPostgresDocumentRepo(p)),
		app.WithAuthService(app.NewCognitoService(c)),
	)
	if err := a.Run(); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
