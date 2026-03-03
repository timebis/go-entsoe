package entsoe

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(
	zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	},
).With().Timestamp().Logger()
