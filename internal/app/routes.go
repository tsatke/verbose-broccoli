package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func (a *App) setupRoutes() {
	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	a.router.Use(sessions.Sessions("SessionID", store))

	rest := a.router.Group("/rest")
	{
		rest.GET("/healthcheck", a.HandlerHealthcheck())
		doc := rest.Group("/doc")
		{
			doc.GET("/:id/content", a.HandlerGetContent())
			doc.POST("/:id/content", a.HandlerPostContent())

			doc.GET("/:id", a.HandlerGetDocument())
			doc.DELETE("/:id", a.HandlerDeleteDocument())

			doc.GET("", a.HandlerGetDocuments())
			doc.POST("", a.HandlerPostDocument())
		}
		user := rest.Group("/user")
		{
			user.POST("/login", a.HandlerUserLogin())
		}
	}
}
