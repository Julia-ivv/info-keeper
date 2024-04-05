package authorizer

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const AccessToken = "accessToken"
const TokenExp = time.Hour * 10

type key string

const UserContextKey key = "token"

type Claims struct {
	jwt.RegisteredClaims
	Login string
	Pwd   string
}

func BuildToken(userLogin, userPwd, secretKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			},
			Login: userLogin,
			Pwd:   userPwd,
		})
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserDataFromToken(tokenString, secretKey string) (userLogin string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if (err != nil) || (!token.Valid) {
		return "", err
	}
	return claims.Login, nil
}
