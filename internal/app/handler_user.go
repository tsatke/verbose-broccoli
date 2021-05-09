package app

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (a *App) HandlerUserLogin() gin.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var r request
		if err := c.BindJSON(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid payload",
			})
			return
		}

		if valid, err := a.users.CredentialsValid(r.Username, r.Password); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Success: false,
				Message: "invalid credentials",
			})
			return
		}

		userID, err := a.users.UserID(r.Username)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		sess := sessions.Default(c)
		sess.Set(UserIDKey, userID)

		if err := sess.Save(); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, Response{
			Success: true,
		})
	}
}
