package log

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultLogger - thread-safe, shared logger.
var DefaultLogger = log.Logger

// Initialises shared thread-safe instanse.
func init() {
	logLevel, ok := os.LookupEnv("POSTMASTER_LOG_LEVEL")
	if ok {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			panic(err)
		}

		zerolog.SetGlobalLevel(level)

		return
	}

	env, ok := os.LookupEnv("POSTMASTER_ENV")
	if strings.EqualFold("production", env) {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		DefaultLogger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// Debug logs message with debug level.
func Debug() *zerolog.Event {
	return DefaultLogger.Debug()
}

// Info logs message with info level.
func Info() *zerolog.Event {
	return DefaultLogger.Info()
}

// Warn logs message with warn level.
func Warn() *zerolog.Event {
	return DefaultLogger.Warn()
}

// Error logs message with error level.
func Error() *zerolog.Event {
	return DefaultLogger.Error()
}

// Fatal logs message with fatal level.
func Fatal() *zerolog.Event {
	return DefaultLogger.Fatal()
}

// Panic logs message with panic level.
func Panic() *zerolog.Event {
	return DefaultLogger.Panic()
}
