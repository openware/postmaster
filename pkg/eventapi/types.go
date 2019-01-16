package eventapi

import "github.com/dgrijalva/jwt-go"

type Event struct {
	Record  map[string]interface{} `json:"record"`
	Changes map[string]interface{} `json:"changes"`
	Name    string                 `json:"name"`
}

type Claims struct {
	jwt.StandardClaims
	Event Event `json:"event"`
}

type DeliverySignatureHeader struct {
	Kid string `json:"kid,omitempty"`
}

type DeliverySignature struct {
	Protected string                  `json:"protected"`
	Signature string                  `json:"signature"`
	Header    DeliverySignatureHeader `json:"header,omitempty"`
}

// Structure of Event API Message.
type Delivery struct {
	Payload    string              `json:"payload"`
	Signatures []DeliverySignature `json:"signatures"`
}
