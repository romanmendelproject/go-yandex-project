package jwt

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v5"
)

const SecretKey = "not secret"

type JWT struct {
}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) GenerateToken(userID int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		log.Error("error generating token", "err", err)

		return "", err
	}

	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SecretKey), nil
	})
	if err != nil {
		log.Error("error while parsing token", "error", err)

		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, errors.New("token is expired")
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, errors.New("invalid token")
		}

		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		log.Errorf("variable type conversion error")

		return 0, fmt.Errorf("variable type conversion error")
	}

	return int(userID), nil
}
