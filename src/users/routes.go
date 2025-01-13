package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/users")

	group.GET("/:id", getUserByIDController)
	group.GET("/verify", verifyEmailController)
	group.POST("/", createUserController)
	group.PUT("/:id/watched/:movie_id", addMovieToUserController)
	group.PUT("/", updateUserController)
	group.DELETE("/:id", deleteUserController)
}
