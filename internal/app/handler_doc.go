package app

import (
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (a *App) HandlerGetContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		if !a.canRead(userID, id) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: "missing read permission",
			})
			return
		}

		content, err := a.objects.Read(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Message: "no content for id",
			})
			return
		}

		_, err = io.Copy(c.Writer, content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Message: "unable to write response",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Success: true,
		})
	}
}

func (a *App) HandlerPostContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		if !a.canWrite(userID, id) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: "missing write permission",
			})
			return
		}

		if err := a.objects.Create(id, c.Request.Body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "failed to create object",
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Success: true,
		})
	}
}

func (a *App) HandlerGetDocument() gin.HandlerFunc {
	type response struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}

	return func(c *gin.Context) {
		id := c.Param("id")
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		if !a.canRead(userID, id) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: "missing read permission",
			})
			return
		}

		header, err := a.index.GetByID(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, Response{
				Message: "failed to obtain document",
			})
			return
		}

		c.JSON(http.StatusOK, response{
			Name: header.Name,
			Size: header.Size,
		})
	}
}

func (a *App) HandlerDeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		if !a.canDelete(userID, id) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: "missing delete permission",
			})
			return
		}

		if err := a.permissions.Delete(userID, id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "failed to delete permission",
			})
			return
		}

		if err := a.index.Delete(id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "failed to delete index entry",
			})
			return
		}

		if err := a.objects.Delete(id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "failed to delete object",
			})
			return
		}
	}
}

func (a *App) HandlerGetDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
	}
}

func (a *App) HandlerPostDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
	}
}
