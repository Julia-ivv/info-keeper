// Пакет authorizer предназначен для авторизации пользователей.
package authorizer

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AccessToken используется для доступа к токену в метаданных gRPC-запроса.
const AccessToken = "accessToken"

// TokenExp - время действия токена.
const TokenExp = time.Hour * 10

type key string

// UserContextKey - для доступа к токену в контексте gRPC-запроса.
const UserContextKey key = "token"

type Claims struct {
	jwt.RegisteredClaims
	Login string
	Pwd   string
}

// BuildToken - создает новый токен.
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

// GetUserDataFromToken - получает логин из токена.
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
