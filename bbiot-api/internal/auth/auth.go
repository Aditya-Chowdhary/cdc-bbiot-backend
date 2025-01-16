package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJWTToken(username string, role string, jwtKey string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["role"] = role

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func getJWTClaims(tokenString string, jwtKey string) (jwt.MapClaims, error) {
	var err error
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid JWT Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims, nil
	}
	return claims, errors.New("invalid claims")
}

func GetUserFromJWTToken(tokenString string, jwtKey string, dbx *database.Queries) (database.User, error) {
	authorizationString := strings.Split(tokenString, " ")

	if len(authorizationString) != 2 {
		return database.User{}, errors.New("invalid Authorization Header")
	}

	token := authorizationString[1]

	if token == "" {
		return database.User{}, errors.New("not logged in")
	}

	claims, err := getJWTClaims(token, jwtKey)

	var team database.User

	if err != nil {
		return team, err
	}
	user, err := dbx.GetUserByUsername(context.Background(), claims["username"].(string))
	if err != nil {
		return user, err
	}

	return user, nil
}
