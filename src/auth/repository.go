package auth

import (
	"errors"
	"log"
	"strings"

	"github.com/SortexGuy/proyecto-db-cassandra/config"
	"github.com/SortexGuy/proyecto-db-cassandra/src/users"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

func loginRepository(loginData LoginDTO) (users.User, error) {
	session := config.SESSION
	user := users.User{}

	// Busca al usuario por el username
	query := `SELECT user_id, email, name, password FROM app.users WHERE name = ?;`
	err := session.Query(query, loginData.Username).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		log.Println("Invalid username: ", err)
		return user, err
	}

	// Verifica el password usando bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		log.Println("Invalid password")
		return users.User{}, errors.New("invalid credentials")
	}

	user.Password = ""
	return user, nil
}

func registrationRepository(registrationData users.User) (users.User, error) {
	session := config.SESSION
	user := users.User{}

	query := `SELECT user_id, name, email FROM app.users`
	iter := session.Query(query).Iter()
	for iter.Scan(&user.ID, &user.Name, &user.Email) {
		registrationData.Email = strings.ToLower(registrationData.Email)
		if user.Email == registrationData.Email || user.Name == registrationData.Name {
			log.Println("name or email already exist")
			return users.User{}, errors.New("name or email already exist")
		}
	}
	if err := iter.Close(); err != nil && err != gocql.ErrNotFound {
		log.Println("Unexpected error: ", err)
		return user, err
	}

	user = registrationData
	query = `INSERT INTO app.users (user_id, name, email, password) VALUES (?, ?, ?, ?)`
	err := session.Query(query, user.ID, user.Name, user.Email, user.Password).Exec()
	if err != nil {
		log.Println("Error inserting user:", err)
		return users.User{}, err
	}

	return user, nil
}
