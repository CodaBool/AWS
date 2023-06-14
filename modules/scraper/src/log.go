package main

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func buildLogger(local bool) {
	if os.Getenv("PRETTY") != "" || !local {
		o := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
		log = zerolog.New(o).With().Logger()
	} else {
		log = zerolog.New(os.Stderr).With().Logger()
		// log = zerolog.New(os.Stderr).With().Str("func", "junk in the trunk").Logger()
	}
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
