package mailconsumer

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type EventAPIClaims struct {
	jwt.StandardClaims

	Event struct {
		Record struct {
			UID                string    `json:"uid"`
			Email              string    `json:"email"`
			Level              int       `json:"level"`
			OtpEnabled         bool      `json:"otp_enabled"`
			ConfirmationToken  string    `json:"confirmation_token"`
			ConfirmationSentAt time.Time `json:"confirmation_sent_at"`
			State              string    `json:"state"`
			FailedAttempts     int       `json:"failed_attempts"`
			CreatedAt          time.Time `json:"created_at"`
			UpdatedAt          time.Time `json:"updated_at"`
		} `json:"record"`
		Name string `json:"name"`
	} `json:"event"`
}

type EventMsgSignature struct {
	Protected string `json:"protected"`
	Header    struct {
		Kid string `json:"kid"`
	} `json:"header"`
	Signature string `json:"signature"`
}

// Structure of Event API Message.
type EventMsg struct {
	Payload    string              `json:"payload"`
	Signatures []EventMsgSignature `json:"signatures"`
}
