package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// MustGetEnv returns error, if variable was not set.
// Otherwise returns content.
func MustGetEnv(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		log.WithFields(log.Fields{"env": name}).Panicln("Environment variable does not set")
	}

	return value
}

// GetEnv returns environment variable with ability to specify default value.
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}
