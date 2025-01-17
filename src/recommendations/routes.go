package recommendations

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/recommendations")
	group.GET("/:user_id", getRecommendationController)
	group.POST("/add", makeRecommendationController)
}
