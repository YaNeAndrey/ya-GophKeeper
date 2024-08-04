package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"ya-GophKeeper/internal/server/secret"
)

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func BuildJWTStringWithLogin(login string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Login:            login,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secret.SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func CheckAccess(r *http.Request) (*Claims, bool) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, false
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret.SecretKey), nil
		})

	if err != nil {
		return nil, false
	}
	if token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}
