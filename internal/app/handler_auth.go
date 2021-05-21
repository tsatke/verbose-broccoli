package app

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (a *App) HandlerAuthLogin() gin.HandlerFunc {
	type request struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	type response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message,omitempty"`
		Challenge string `json:"challenge,omitempty"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid JSON payload",
			})
			return
		}

		res, err := a.auth.Login(req.Username, req.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
			})
			return
		}

		if !res.Success {
			c.JSON(http.StatusUnauthorized, response{
				Success: false,
				Message: "invalid credentials",
			})
			return
		}

		if !res.Done {
			c.JSON(http.StatusOK, response{
				Success:   true,
				Message:   "request not done yet",
				Challenge: res.Challenge,
			})
			return
		}

		sess := sessions.Default(c)
		sess.Set(UserIDKey, res.Token)
		if err := sess.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Success: false,
				Message: "unable to save session",
			})
			return
		}

		c.JSON(http.StatusOK, response{
			Success: true,
		})
	}
}

func (a *App) HandlerAuthChallenge() gin.HandlerFunc {
	type request struct {
		Challenge      string `json:"challenge,omitempty"`
		ClientResponse string `json:"client_response,omitempty"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid JSON payload",
			})
			return
		}
	}
}

func (a *App) HandlerAuthLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Clear()

		if err := sess.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Success: false,
				Message: "unable to save session",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Success: true,
		})
	}
}
