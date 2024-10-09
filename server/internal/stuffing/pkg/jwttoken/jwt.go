package jwttoken

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/golang-jwt/jwt/v4"
)

var refreshSecret = []byte("your_refresh_secret_key")
var jwtKey = []byte("Fdhc53$537&20dkjfG")

type Claims struct {
	User string //user.User
	jwt.RegisteredClaims
}

type SetJWTToken struct{}

func SetToken(user *user.User) (string, error) {
	user.ID_User = user.ID.Hex()
	expirationTimeSecond := 20
	userByte, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error marshalling user into set token")
	}

	expirationTime := time.Now().Add(time.Duration(expirationTimeSecond) * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		User: string(userByte),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}

	return tokenString, nil
}

// Функция для создания Refresh Token
func CreateRefreshToken(user *user.User) (string, error) {
	user.ID_User = user.ID.Hex()
	userByte, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error marshalling user into set token")
	}
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh Token действителен 7 дней
		Subject:   string(userByte),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// Функция для обновления JWT с помощью Refresh Token
func RefreshJWT(refreshToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	user := user.User{}

	err = json.Unmarshal([]byte(claims.Subject), &user)

	if err != nil {
		fmt.Println("error unmarsh claims user")
		return "", err
	}

	return SetToken(&user) // Здесь можно получить роль пользователя из базы данных или другого источника
}

func Welcom(token string) (*user.User, error) {
	// Get the JWT string from the cookie
	tknStr := token

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if tkn != nil {
		if !tkn.Valid {
			return nil, err
		}
	} else {
		return nil, err
	}
	user := user.User{}

	err = json.Unmarshal([]byte(claims.User), &user)

	if err != nil {
		fmt.Println("error unmarsh claims user")
		return nil, err
	}

	return &user, nil
}
