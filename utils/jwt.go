package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

package utils

import (
"errors"
"github.com/dgrijalva/jwt-go"
)

// 用于签发和验证JWT
const secretKey = "bz2021ABCDefgHIJKlmn987QMWNEB11"

// GenerateJWT 生成JWT
func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT 验证和解析JWT
func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("invalid user ID")
	}

	return userID, nil
}

