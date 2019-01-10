package eventapi

import (
	"encoding/base64"
	"errors"
)

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/shal/pigeon/pkg/utils"
)

func ValidateJWT(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("unexpected signing method")
	}

	encPublicKey := utils.MustGetEnv("JWT_PUBLIC_KEY")
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
