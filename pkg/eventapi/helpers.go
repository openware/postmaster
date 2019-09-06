package eventapi

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Unmarshal returns Event from RawEvent.
func Unmarshal(raw RawEvent) (*Event, error) {
	event := Event{}

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

	return &event, nil
}

// FixAndValidate returns error in case of broken event payload.
func (event *Event) FixAndValidate(language string) (*Record, error) {
	var record *Record

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           &record,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, fmt.Errorf("mapstructure: %s", err)
	}

	if err := decoder.Decode(event.Record); err != nil {
		return nil, fmt.Errorf("decorer: %s", err)
	}

	if record == nil {
		return nil, errors.New("event: record is nil")
	}

	if record.User == nil {
		return nil, errors.New("event: record.user is nil")
	}

	if record.User.UID == "" {
		return nil, errors.New("event: record.user.uid is empty")
	}

	if record.User.Email == "" {
		return nil, errors.New("event: record.user.email is empty")
	}

	// First language in config is default.
	if record.Language == "" {
		record.Language = language
	}

	return record, nil // Everythin is valid.
}
