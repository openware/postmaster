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

type Record struct {
	User     User   `json:"user"`
	Language string `json:"language"`
}

type Event struct {
	Record  map[string]interface{} `json:"record"`
	Changes map[string]interface{} `json:"changes"`
}

func Unmarshal(raw RawEvent) (*Event, error) {
	event := new(Event)
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &event,
		WeaklyTypedInput: true,
	})

	if err != nil {
		return nil, err
	}

	if err := dec.Decode(raw); err != nil {
		return nil, err
	}

	return event, nil
}
