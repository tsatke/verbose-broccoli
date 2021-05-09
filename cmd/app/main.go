package main

import (
	"github.com/tsatke/verbose-broccoli/internal/app"
	"github.com/tsatke/verbose-broccoli/internal/app/mem"
)

func main() {
	u := mem.NewUserService()
	_ = u.CreateUser("foo", "bar")

	a := app.New(nil, u, nil, nil, nil)
	_ = a.Run(":8080")
}
