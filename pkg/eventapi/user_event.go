package eventapi

import (
	"github.com/mitchellh/mapstructure"
)

type User struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Level int    `json:"level"`
	Otp   bool   `json:"otp_enabled"`
	State string `json:"state"`
}

type UserEvent struct {
	User     User   `json:"user"`
	Language string `json:"language"`
}

func Unmarshal(event Event) (*UserEvent, error) {
	usrEvent := new(UserEvent)
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &usrEvent,
		WeaklyTypedInput: true,
	})

	if err != nil {
		return nil, err
	}

	if err := dec.Decode(event); err != nil {
		return nil, err
	}

	return usrEvent, nil
}
