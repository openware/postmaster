package eventapi

import (
	"encoding/base64"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type RawEvent map[string]interface{}

type Claims struct {
	jwt.StandardClaims
	Event RawEvent `json:"event"`
}

type Validator struct {
	Algorithm string `yaml:"algorithm"`
	Value     string `yaml:"value"`
}

func (v *Validator) ValidateJWT(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != v.Algorithm {
		return nil, errors.New("unexpected signing method")
	}

	publicKey, err := base64.StdEncoding.DecodeString(v.Value)

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
