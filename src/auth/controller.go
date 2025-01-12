package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginUserController(c *gin.Context) {
	loginDTO := LoginDTO{}
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//token, err := loginService(loginDTO.Username, loginDTO.Password)
	token, err := loginService(loginDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func registerUserController(c *gin.Context) {
	registrationDTO := RegistrationDTO{}

	if err := c.ShouldBindJSON(&registrationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := registrationService(registrationDTO)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
