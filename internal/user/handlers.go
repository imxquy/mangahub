package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mangahub/internal/auth"
)

func RegisterLibraryRoutes(authed *gin.RouterGroup, repo *Repo) {
	authed.POST("/users/library", func(c *gin.Context) {
		var req AddLibraryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		userID := auth.GetUserID(c)
		if err := repo.AddToLibrary(userID, req.MangaID, req.Status); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"ok": true})
	})

	authed.GET("/users/library", func(c *gin.Context) {
		userID := auth.GetUserID(c)
		items, err := repo.ListLibrary(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})
}
