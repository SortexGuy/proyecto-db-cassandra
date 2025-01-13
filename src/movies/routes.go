package movies

import "github.com/gin-gonic/gin"

// RegisterRoutes registra las rutas para el grupo de películas
func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/movies")

	// CRUD básico
	group.GET("/:id", getMovieByIDController)
	group.POST("/", createMovieController)
	group.PUT("/", updateMovieController)
	group.DELETE("/:id", deleteMovieController)
	group.GET("/", getAllMoviesController)

	// Relación usuario-película
	group.GET("/user/:user_id", getMovieWatchedByUserController)
}
