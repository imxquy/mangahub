package manga

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repo *Repo) {
	r.GET("/manga", func(c *gin.Context) {
		q := c.Query("q")
		genre := c.Query("genre")
		status := c.Query("status")
		limit := 20

		results, err := repo.Search(q, genre, status, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	})

	r.GET("/manga/:id", func(c *gin.Context) {
		id := c.Param("id")
		m, err := repo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "manga not found"})
			return
		}
		c.JSON(http.StatusOK, m)
	})
}
