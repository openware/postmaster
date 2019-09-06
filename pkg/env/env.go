package env

import (
	"fmt"
	"os"

	"github.com/openware/postmaster/internal/log"
)

func Must(value string, err error) string {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return value
}

// Fetch returns value of environment variable.
// Error returned in case of empty environment variable value.
func Fetch(name string) (string, error) {
	value, exist := os.LookupEnv(name)

	if !exist {
		return "", fmt.Errorf("environment variable %s does not set", name)
	}

	return value, nil
}

// FetchDefault returns environment variable with ability to specify default value.
func FetchDefault(key, fallback string) string {
	if value, err := Fetch(key); err == nil {
		return value
	}

	return fallback
}
