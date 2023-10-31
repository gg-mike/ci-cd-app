package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func MultiOutput(logDir string) {
	os.MkdirAll(logDir, os.ModePerm)
	f, err := os.Create(fmt.Sprintf("%s/server-%s.log", logDir, time.Now().Format("20060102-150405")))
	if err != nil {
		log.Fatal().Msgf("Error creating log file (%v)", err)
	}
	zerolog.NewConsoleWriter()
	writers := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}, f)
	log.Logger = log.Output(writers)
}
