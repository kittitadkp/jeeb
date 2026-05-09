package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	w := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}
	log = zerolog.New(w).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// Init configures the global logger. Call once from PersistentPreRunE.
// level is one of: debug, info, warn, error (case-insensitive; invalid values default to info).
func Init(level string) {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	w := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}
	log = zerolog.New(w).With().Timestamp().Logger()
}

// Step writes a user-facing progress line to stdout. Always visible regardless of log level.
func Step(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// StepMsg writes a plain user-facing line to stdout.
func StepMsg(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

func Debug(msg string, args ...any) {
	if len(args) > 0 {
		log.Debug().Msgf(msg, args...)
	} else {
		log.Debug().Msg(msg)
	}
}

func Info(msg string, args ...any) {
	if len(args) > 0 {
		log.Info().Msgf(msg, args...)
	} else {
		log.Info().Msg(msg)
	}
}

func Warn(msg string, args ...any) {
	if len(args) > 0 {
		log.Warn().Msgf(msg, args...)
	} else {
		log.Warn().Msg(msg)
	}
}

func Error(msg string, args ...any) {
	if len(args) > 0 {
		log.Error().Msgf(msg, args...)
	} else {
		log.Error().Msg(msg)
	}
}
