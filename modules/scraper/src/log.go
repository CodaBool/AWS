package main

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func buildLogger() {
	logger = zerolog.New(os.Stderr).With().Logger()
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		o := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
		logger = zerolog.New(o).With().Logger()
	}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	if os.Getenv("QUIET") != "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func check(err error, log ...zerolog.Logger) {
	// closest thing go has to default arguments
	if len(log) > 0 {
		if err != nil {
			log[0].Fatal().Err(err).Msg("")
		}
	} else {
		if err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
}
