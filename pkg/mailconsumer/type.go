package mailconsumer

import "time"

type JWTPayload struct {
	Iss   string `json:"iss"`
	Jti   string `json:"jti"`
	Iat   int    `json:"iat"`
	Exp   int    `json:"exp"`
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
	Alg string `json:"alg"`
}