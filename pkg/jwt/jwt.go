package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func ExtractClaims(tokenStr string, signingKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}
