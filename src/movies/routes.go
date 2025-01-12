package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas para el grupo de películas
func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/movies")

	group.GET("/", getMovieByIDController)
	group.GET("/", getMovieByUserController)
	group.GET("/", getAllMoviesController)
}
