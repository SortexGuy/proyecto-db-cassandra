package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var fallback = []byte("tu_clave_secreta")

func loginService(loginData LoginDTO) (string, error) {
	_, err := loginRepository(loginData)
	if err != nil {
		return "", err
	}

	if jwtKey == nil {
		jwtKey = fallback
	}

	expirationTime := time.Now().Add(24 * time.Hour) // 24h
	var claims = Claims{
		Username: loginData.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func registrationService(registrationData RegistrationDTO) (string, error) {
	_, err := registrationRepository(registrationData)
	if err != nil {
		return "", err
	}

	if jwtKey == nil {
		jwtKey = fallback
	}

	expirationTime := time.Now().Add(24 * time.Hour) // 24h
	var claims = Claims{
		Username: registrationData.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
