package auth

import (
	"github.com/SortexGuy/proyecto-db-cassandra/src/counters"
	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
	"golang.org/x/crypto/bcrypt"
)

func loginService(loginData LoginDTO) (users.User, error) {
	user, err := loginRepository(loginData)

	return user, err
}

func registrationService(registrationData RegistrationDTO) (users.User, error) {
	user := users.User{}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	id, err := counters.IncrementCounter("users")
	if err != nil {
		return user, err
	}
	user.ID = id
	user.Password = string(hashedPassword)

	user, err = registrationRepository(user)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}
