package eventapi

import (
	"encoding/base64"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// RawEvent is raw JSON payload of incoming event.
type RawEvent map[string]interface{}

// Claims is JWT Token claims.
type Claims struct {
	jwt.StandardClaims
	Event RawEvent `json:"event"`
}

// Validator is JSON Web Token validator.
type Validator struct {
	Algorithm string `yaml:"algorithm"`
	Value     string `yaml:"value"`
}

// ValidateJWT validates, that JWT token is properly signed.
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

// ParseJWT parses JSON Web Token and returns ready for use claims.
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
