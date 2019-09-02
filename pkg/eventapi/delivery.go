package eventapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/streadway/amqp"
)

type DeliverySignatureHeader struct {
	Kid string `json:"kid,omitempty"`
}

type DeliverySignature struct {
	Protected string                  `json:"protected"`
	Signature string                  `json:"signature"`
	Header    DeliverySignatureHeader `json:"header,omitempty"`
}

type Delivery struct {
	Payload    string              `json:"payload"`
	Signatures []DeliverySignature `json:"signatures"`
}

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
