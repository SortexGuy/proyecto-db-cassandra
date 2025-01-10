package users

import (
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
)

func getUsersController(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id not found",
		})
		return
	}

	user, err := GetUserByIDService(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting user",
		})
		return
	}

	data := []UserDTO{user}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}


// UserController es la estructura que maneja las operaciones relacionadas con los usuarios
type UserController struct {
	userRepo *UserRepository
}

// NewUser Controller crea una nueva instancia de UserController
func NewUserController(userRepo *UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

// GetUsers maneja la solicitud para obtener todos los usuarios
func (ctrl *UserController) GetUsers() ([]User , error) {
	users, err := ctrl.userRepo.GetAllUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}

	// Contabilizar cuántos usuarios se extrajeron
	count := len(users)
	log.Printf("Total users extracted: %d\n", count)

	return users, nil
}