package eventapi

import "github.com/dgrijalva/jwt-go"

type Event map[string]interface{}

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

type User struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Level int    `json:"level"`
	Otp   bool   `json:"otp_enabled"`
	State string `json:"state"`
}
