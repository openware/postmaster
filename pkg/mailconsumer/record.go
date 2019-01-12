package mailconsumer

import (
	"fmt"
)

import (
	"github.com/shal/pigeon/pkg/utils"
)

type AccountRecord struct {
	ConfirmationToken string `json:"confirmation_token"`
	UID               string `json:"uid"`
	Email             string `json:"email"`
	Level             int    `json:"level"`
	OtpEnabled        bool   `json:"otp_enabled"`
	State             string `json:"state"`
	FailedAttempts    int    `json:"failed_attempts"`
}

func (record AccountRecord) ConfirmationUri() string {
	base := utils.GetEnv("FRONTEND_DOMAIN", "http://www.example.com")

	return fmt.Sprintf("%s/accounts/confirmation?confirmation_token=%s",
		base,
		record.ConfirmationToken,
	)
}
