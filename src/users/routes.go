package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/users")

	group.GET("/", getUserByIDController)
	group.GET("/verify", verifyEmailController)
	group.POST("/", insertUserController)
	group.PUT("/", updateUserController)
	group.PUT("/watched", updateWatchedMovieController)
	group.DELETE("/", deleteUserController)
}
