package movies

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/movies")

	group.GET("/", getMoviesController)
}
