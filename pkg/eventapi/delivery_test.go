package eventapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestOneSigDeliveryAsJWT(t *testing.T) {
	delivery := Delivery{
		Payload: "y",
		Signatures: []DeliverySignature{
			{
				Protected: "x",
				Signature: "z",
			},
		},
	}

	body, _ := json.Marshal(delivery)
	res, err := DeliveryAsJWT(amqp.Delivery{Body: body})

	assert.NoError(t, err)
	assert.Equal(t, strings.NewReader("x.y.z"), res)
}

func TestMultiSigDeliveryAsJWT(t *testing.T) {
	delivery := Delivery{
		Payload: "y",
		Signatures: []DeliverySignature{
			{
				Protected: "x",
				Signature: "z",
			},
			{
				Protected: "x",
				Signature: "z",
			},
		},
	}

	body, _ := json.Marshal(delivery)
	res, err := DeliveryAsJWT(amqp.Delivery{Body: body})

	assert.Nil(t, res)
	assert.Equal(t, "multi signature JWT keys does not supported", err.Error())
}

func TestEmptyDeliveryAsJWT(t *testing.T) {
	res, err := DeliveryAsJWT(amqp.Delivery{Body: []byte("{}")})

	assert.Nil(t, res)
	assert.Equal(t, "no signatures to verify", err.Error())
}
