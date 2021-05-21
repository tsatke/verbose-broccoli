package app

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (a *App) setupRoutes() {
	a.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	a.router.Use(sessions.Sessions("SessionID", store))

	a.router.Use(func(c *gin.Context) {
		// don't check the token for the these routes
		switch c.FullPath() {
		case "/rest/healthcheck",
			"/rest/auth/login":
			return
		}

		sess := sessions.Default(c)
		if sess.Get(UserIDKey) == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "not logged in",
			})
			return
		}
	})

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
		auth := rest.Group("/auth")
		{
			auth.POST("/login", a.HandlerAuthLogin())
			auth.GET("/logout", a.HandlerAuthLogout())
			auth.POST("/challenge", a.HandlerAuthChallenge())
		}
	}
}
