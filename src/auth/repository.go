package auth

// import (
// 	"errors"
//
// 	"github.com/SortexGuy/proyecto-db-cassandra/config"
// )

func loginRepository(loginData LoginDTO) (LoginDTO, error) {
	// db := config.SESSION
	//
	// // Busca al usuario por el username
	// if err := db.Where("username = ?", loginData.username).First(&user).Error; err != nil {
	// 	return nil, errors.New("invalid credentials")
	// }
	//
	// // Verifica el password usando bcrypt
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return nil, errors.New("invalid credentials")
	// }

	return loginData, nil
}

func registrationRepository(registrationData RegistrationDTO) (RegistrationDTO, error) {

	return registrationData, nil
}
