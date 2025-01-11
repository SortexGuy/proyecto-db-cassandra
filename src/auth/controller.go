package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginController(c *gin.Context) {
	loginDTO := LoginDTO{}
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := loginService(loginDTO.Username, loginDTO.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
