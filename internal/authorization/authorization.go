package authorization

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/encrypt"
	"github.com/poggerr/gophermart/internal/logger"
	"github.com/poggerr/gophermart/internal/models"
	"github.com/poggerr/gophermart/internal/storage"
	"os"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID *uuid.UUID
}

const TokenExp = time.Hour * 3

func BuildJWTString(uuid *uuid.UUID) (string, error) {

	var secretKey = os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: uuid,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserID(tokenString string) *uuid.UUID {
	//var secretKey = os.Getenv("SECRET_KEY")
	var secretKey = "scdcsdc,HVJHVCAJscdJccdsJVDVJDvqwe[p[;cqsc09cah989h"
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return claims.UserID
}

func RegisterUser(strg *storage.Storage, user *models.User) (uuid.UUID, error) {
	user.Password = encrypt.Encrypt(user.Password)
	id := uuid.New()
	err := strg.CreateUser(user.Username, user.Password, &id)
	if err != nil {
		logger.Initialize().Error(err)
		return id, err
	}
	return id, nil
}
