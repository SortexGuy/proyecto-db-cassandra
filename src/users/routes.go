package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/users")

	group.GET("/", getUserByIDController)
}
