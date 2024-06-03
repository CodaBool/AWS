package main

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func buildLogger() {
	log = zerolog.New(os.Stderr).With().Logger()
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		o := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
		log = zerolog.New(o).With().Logger()
	}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	if os.Getenv("QUIET") != "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
