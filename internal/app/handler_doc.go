package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (a *App) HandlerGetContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := DocID(c.Param("id"))
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

		content, err := a.objects.Read(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Message: "no content for id",
			})
			return
		}

		_, err = io.Copy(c.Writer, content)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
}

func (a *App) HandlerPostContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := DocID(c.Param("id"))
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		acl, err := a.documents.ACL(id)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "get ACL for document",
			})
			return
		}
		if _, ok := acl.Permissions[userID]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: "write content",
			})
			return
		}

		ff, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Message: "failed to receive file",
			})
			return
		}

		f, err := ff.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Message: "failed to open file",
			})
			return
		}
		defer func() {
			_ = f.Close()
		}()

		if err := a.objects.Create(id, f); err != nil {
			_ = c.Error(err)
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
	}

	return func(c *gin.Context) {
		id := DocID(c.Param("id"))
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

		header, err := a.documents.Get(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, Response{
				Message: "failed to obtain document",
			})
			return
		}

		c.JSON(http.StatusOK, response{
			Name: header.Name,
		})
	}
}

func (a *App) HandlerDeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := DocID(c.Param("id"))
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

		// if err := a.permissions.Delete(userID, id); err != nil {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		// 		Message: "failed to delete permission",
		// 	})
		// 	return
		// }

		if err := a.documents.Delete(id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Message: "failed to delete documents entry",
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
	type request struct {
		Filename string `json:"filename"`
	}
	type response struct {
		Success bool   `json:"success"`
		ID      DocID  `json:"id,omitempty"`
		Message string `json:"message,omitempty"`
	}
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		userID := sess.Get(UserIDKey).(string)

		var req request
		if err := c.ShouldBindJSON(&req); err != nil || req.Filename == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid JSON payload",
			})
			return
		}

		id := DocID(a.genUUID().String())

		if err := a.documents.Create(DocumentHeader{
			ID:      id,
			Name:    req.Filename,
			Owner:   userID,
			Created: a.clock.Now(),
		}, ACL{
			Permissions: map[string]Permission{
				userID: {
					Username: userID,
					Read:     true,
					Write:    true,
					Delete:   true,
					Share:    true,
				},
			},
		}); err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Success: false,
				Message: "unable to create documents entry",
			})
			return
		}

		c.JSON(http.StatusOK, response{
			Success: true,
			ID:      id,
		})
	}
}
