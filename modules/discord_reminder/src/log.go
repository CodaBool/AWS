package main

import (
	"fmt"
	log "log/slog"
	"os"
	"strconv"
)

var slog *log.Logger

func buildLogger(debug bool, verbose bool, local bool) {
	var RESET = "\033[0m"
	var RED_BG = "\033[41m"
	var ORANGE_BG = "\033[43m"
	var GREEN_BG = "\033[42m"
	var PURPLE_BG = "\033[45m"

	level := log.LevelInfo
	if debug {
		level = log.LevelDebug
	}

	config := &log.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a log.Attr) log.Attr {
			if a.Key == log.TimeKey {
				a.Value = log.StringValue(a.Value.Time().Format("01/02 15:04:05"))
				if !verbose {
					return log.Attr{}
				}
			} else if a.Key == log.SourceKey {
				source := a.Value.Any().(*log.Source)
				a.Value = log.StringValue(source.Function + ":" + strconv.Itoa(source.Line))
				a.Key = "src"
				if !verbose {
					return log.Attr{}
				}
			} else if a.Key == log.LevelKey {
				if !local {
					return a
				}
				level := a.Value.Any().(log.Level)
				if level == log.LevelDebug {
					fmt.Print(PURPLE_BG + " " + RESET + " ")
				} else if level == log.LevelInfo {
					fmt.Print(GREEN_BG + " " + RESET + " ")
				} else if level == log.LevelWarn {
					fmt.Print(ORANGE_BG + " " + RESET + " ")
				} else {
					fmt.Print(RED_BG + " " + RESET + " ")
				}
				return log.Attr{}
			} else if a.Key == log.MessageKey {
				if local {
					fmt.Print(a.Value.Any().(string) + " ")
					return log.Attr{}
				}
			}
			return a
		},
	}
	if local {
		slog = log.New(log.NewTextHandler(os.Stdout, config))
	} else {
		slog = log.New(log.NewJSONHandler(os.Stdout, config))
	}
	log.SetDefault(slog)
}

func check(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}
