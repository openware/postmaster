package mailconsumer

import (
	"encoding/json"
	"testing"
)

import (
	"github.com/shal/pigeon/pkg/eventapi"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestOneSigDeliveryAsJWT(t *testing.T) {
	delivery := eventapi.Delivery{
		Payload: "y",
		Signatures: []eventapi.DeliverySignature{
			{
				Protected: "x",
				Signature: "z",
			},
		},
	}

	body, _ := json.Marshal(delivery)
	res, err := DeliveryAsJWT(amqp.Delivery{ Body: body })

	assert.NoError(t, err)
	assert.Equal(t, "x.y.z", res)
}

func TestMultiSigDeliveryAsJWT(t *testing.T) {
	delivery := eventapi.Delivery{
		Payload: "y",
		Signatures: []eventapi.DeliverySignature{
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
	_, err := DeliveryAsJWT(amqp.Delivery{ Body: body })

	assert.Equal(t, "multi signature JWT keys does not supported", err.Error())
}

func TestDeliveryAsJWT(t *testing.T) {
	_, err := DeliveryAsJWT(amqp.Delivery{ Body: []byte("{}") })
	assert.Equal(t, "no signatures to verify", err.Error())
}

