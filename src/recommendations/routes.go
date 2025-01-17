package recommendations

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/recommendations")
	group.POST("/add", makeRecommendationController)
	group.GET("/", getRecommendationController)
}
