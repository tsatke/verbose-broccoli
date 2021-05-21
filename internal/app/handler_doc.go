package app

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *App) HandlerGetContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
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
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

		ff, err := c.FormFile("file")
		if err != nil {
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
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

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
		// sess := sessions.Default(c)
		// userID := sess.Get(UserIDKey).(string)

		// if err := a.permissions.Delete(userID, id); err != nil {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		// 		Message: "failed to delete permission",
		// 	})
		// 	return
		// }

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
	type request struct {
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}
	type response struct {
		Success bool   `json:"success"`
		ID      string `json:"id,omitempty"`
		Message string `json:"message,omitempty"`
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil || req.Filename == "" || req.Size <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid JSON payload",
			})
			return
		}

		id := uuid.New().String()

		if err := a.index.Create(DocumentHeader{
			ID:   id,
			Name: req.Filename,
			Size: req.Size,
		}); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Success: false,
				Message: "unable to create index entry",
			})
			return
		}

		c.JSON(http.StatusOK, response{
			Success: true,
			ID:      id,
		})
	}
}
