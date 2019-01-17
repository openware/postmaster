package eventapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/streadway/amqp"
	"io"
	"strings"
)

func DeliveryAsJWT(delivery amqp.Delivery) (io.Reader, error) {
	eventMsg := Delivery{}

	if err := json.Unmarshal(delivery.Body, &eventMsg); err != nil {
		return nil, err
	}

	// Verify JWT payload.
	if len(eventMsg.Signatures) < 1 {
		return nil, errors.New("no signatures to verify")
	} else if len(eventMsg.Signatures) > 1 {
		return nil, errors.New("multi signature JWT keys does not supported")
	}

	// Build token from received header, payload, signatures.
	tokenStr := fmt.Sprintf("%s.%s.%s",
		eventMsg.Signatures[0].Protected,
		eventMsg.Payload,
		eventMsg.Signatures[0].Signature,
	)

	return strings.NewReader(tokenStr), nil
}

func ParseJWT(tokenStr string, callback func(map[string]interface{})) error {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, ValidateJWT)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return errors.New("claims: invalid jwt token")
	}

	// Send email.
	callback(claims.Event.Record)

	return nil
}
