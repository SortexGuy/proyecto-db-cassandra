package recommendations

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func makeRecommendationController(c *gin.Context) {
	userIDStr := c.Query("user_id")
	lambdaStr := c.Query("lambda")
	// Validar parámetros
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	lambda, err := strconv.ParseFloat(lambdaStr, 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie_id"})
		return
	}

	log.Println("Entrando al servicio")
	// Llamar al servicio
	r, err := makeRecommendationService(userID, lambda)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed saving data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": r})
}
