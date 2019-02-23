package eventapi

import (
	"encoding/base64"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/openware/postmaster/pkg/env"
)

type Event map[string]interface{}

type Claims struct {
	jwt.StandardClaims
	Event Event `json:"event"`
}

func ValidateJWT(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("unexpected signing method")
	}

	encPublicKey := env.Must(env.Fetch("JWT_PUBLIC_KEY"))
	publicKey, err := base64.StdEncoding.DecodeString(encPublicKey)

	if err != nil {
		return nil, err
	}

	signingKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		return nil, err
	}

	return signingKey, nil
}

func ParseJWT(tokenStr string, keyFunc jwt.Keyfunc) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("claims: invalid jwt token")
	}

	return claims, nil
}
