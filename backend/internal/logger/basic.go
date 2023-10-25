package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Basic(level zerolog.Level, module string) *zerolog.Event {
	if level == zerolog.FatalLevel {
		return log.Logger.Fatal().Str("module", module)
	} else if level == zerolog.PanicLevel {
		return log.Logger.Panic().Str("module", module)
	}
	return log.Logger.WithLevel(level).Str("module", module)
}
