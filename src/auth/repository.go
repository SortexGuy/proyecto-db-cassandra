package auth

import (
// "errors"
//
// "golang.org/x/crypto/bcrypt"
)

func loginRepository(username, password string) (*LoginDTO, error) {
	// db := config.DB
	// var user schema.User

	// Busca al usuario por el username
	// if err := db.Where("username = ?", username).First(&user).Error; err != nil {
	// 	return nil, errors.New("invalid credentials")
	// }
	//
	// // Verifica el password usando bcrypt
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return nil, errors.New("invalid credentials")
	// }
	//
	// return &user, nil
	return &LoginDTO{}, nil
}
