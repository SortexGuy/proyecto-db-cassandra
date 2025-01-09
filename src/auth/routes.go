package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/auth")

	group.GET("/", getAuthController)
}
