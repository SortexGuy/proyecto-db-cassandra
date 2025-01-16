package auth

import (
	"net/http"

	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
	"github.com/gin-gonic/gin"
)

func loginUserController(c *gin.Context) {
	loginDTO := LoginDTO{}
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := loginService(loginDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []users.User{user},
	})
}

func registerUserController(c *gin.Context) {
	registrationDTO := RegistrationDTO{}

	if err := c.ShouldBindJSON(&registrationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := registrationService(registrationDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []users.User{user},
	})
}
