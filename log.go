package main

import (
	"os"

	"github.com/rs/zerolog"
)

func createLogger() zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Level(zerolog.InfoLevel)

	return logger
}
