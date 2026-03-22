package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(level, env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	if env == "development" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(lvl).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return zerolog.New(os.Stderr).
		Level(lvl).
		With().
		Timestamp().
		Logger()
}
