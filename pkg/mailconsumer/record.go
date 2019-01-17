package mailconsumer

import (
	"strings"
)

import (
	"github.com/shal/pigeon/pkg/utils"
)

type AccountRecord struct {
	Email             string `json:"email"`
	ConfirmationToken string `json:"confirmation_token"`
	UID               string `json:"uid"`
	Level             int    `json:"level"`
	OtpEnabled        bool   `json:"otp_enabled"`
	State             string `json:"state"`
	FailedAttempts    int    `json:"failed_attempts"`
}

func (r AccountRecord) ConfirmationUri() string {
	//base := utils.GetEnv("FRONTEND_DOMAIN", "http://www.example.com")
	url := utils.GetEnv("CONFIRM_URL",
		"http://www.example.com/accounts/confirmation?confirmation_token=#{}",
	)

	return strings.Replace(url, "#{}", r.ConfirmationToken, 1)
}
