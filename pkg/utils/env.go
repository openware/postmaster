package utils

import (
	"fmt"
	"os"
)

// MustGetEnv returns error, if variable was not set.
// Otherwise returns content.
func MustGetEnv(name string) string {
	value, exist := os.LookupEnv(name)

	if !exist {
		panic(fmt.Sprintf("Environment variable %s does not set", name))
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
