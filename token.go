package kafekoding_api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type credential struct {
	ID int `json:"id"`
}

type claims struct {
	credential
	jwt.RegisteredClaims
}

type userAuth struct{}

var secretKey = []byte("Secrey Key")

// generateNewToken is function to get new token for user.
func generateNewToken(credential credential) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := claims{
		credential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(secretKey)
	return token, err
}

// verifyToken is function to check the token is valid or not.
func verifyToken(token string) (claims, error) {
	var claims claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return claims, err
	}

	if !jwtToken.Valid {
		return claims, errors.New("token anda sudah tidak valid atau kadaluarsa")
	}

	return claims, nil
}
