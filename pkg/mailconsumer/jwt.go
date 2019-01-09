package mailconsumer

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/shal/mail-consumer/pkg/utils"
)

func ValidateJWT(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("unexpected signing method")
	}

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")

	encPublicKey := utils.MustGetEnv("JWT_PUBLIC_KEY")
	publicKey, err := base64.StdEncoding.DecodeString(encPublicKey)

	fmt.Println(string(publicKey))

	if err != nil {
		return nil, err
	}

	signingKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		return nil, err
	}

	return signingKey, nil
}
