package mailconsumer

// TODO: Ensure, that this structure needed.
type AccountRecord struct {
	UID               string `json:"uid"`
	Email             string `json:"email"`
	Level             int    `json:"level"`
	OtpEnabled        bool   `json:"otp_enabled"`
	ConfirmationToken string `json:"confirmation_token"`
	State             string `json:"state"`
	FailedAttempts    int    `json:"failed_attempts"`

	// ConfirmationSentAt time.Time `json:"confirmation_sent_at"`
	// CreatedAt time.Time `json:"2019-01-10T14:20:36Z"`
	// UpdatedAt time.Time `json:"2019-01-10T14:20:36Z"`
}
