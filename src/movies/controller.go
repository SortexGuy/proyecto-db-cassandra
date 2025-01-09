package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMoviesController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "movies",
	})
}
