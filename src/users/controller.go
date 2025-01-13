package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createUserController(c *gin.Context) {
	var userDTO User
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := createUserService(userDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func getUserByIDController(c *gin.Context) {
	userIDText := c.Param("id")
	userID, err := strconv.ParseInt(userIDText, 10, 64)
	if userIDText == "" || err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "movie_id not found",
		})
		return
	}

	user, err := getUserByIDService(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting user",
		})
		return
	}

	data := []User{user}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func verifyEmailController(c *gin.Context) {
	emailText := c.Query("email")
	if emailText == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "email text not found",
		})
		return
	}

	exist, err := verifyEmailService(emailText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting user",
		})
		return
	}

	data := []bool{exist}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// Controlador para actualizar un usuario
func updateUserController(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := updateUserService(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func addMovieToUserController(c *gin.Context) {
	// Extraer userID y movieID de los parámetros
	userIDParam := c.Param("id")
	movieIDParam := c.Param("movie_id")

	// Convertir parámetros a int64
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	movieID, err := strconv.ParseInt(movieIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Llamar al servicio para registrar la película
	err = addMovieToUserService(userID, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie marked as watched successfully"})
}

// Controlador para eliminar un usuario
func deleteUserController(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = deleteUserService(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
