package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/users")

	group.GET("/", getUserByIDController)
	group.POST("/", insertUserController)
	group.PUT("/", updateUserController)
	group.DELETE("/", deleteUserController)
}
