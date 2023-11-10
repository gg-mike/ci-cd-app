package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Debug(module string) *zerolog.Event {
	return log.Logger.Debug().Str("module", module)
}

func Info(module string) *zerolog.Event {
	return log.Logger.Info().Str("module", module)
}

func Warn(module string) *zerolog.Event {
	return log.Logger.Warn().Str("module", module)
}

func Error(module string) *zerolog.Event {
	return log.Logger.Error().Str("module", module)
}

func Fatal(module string) *zerolog.Event {
	return log.Logger.Fatal().Str("module", module)
}

func Panic(module string) *zerolog.Event {
	return log.Logger.Panic().Str("module", module)
}

func Trace(module string) *zerolog.Event {
	return log.Logger.Trace().Str("module", module)
}
