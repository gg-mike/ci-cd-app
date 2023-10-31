package main

import (
	_ "github.com/gg-mike/ci-cd-app/backend/docs/server"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/sys"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Basic(zerolog.WarnLevel, "main").Msgf("Missing .env file")
	}

	logDir := sys.GetEnvWithFallback("LOGS", "./logs")
	logger.MultiOutput("dbMigrator", logDir)

	dbUrl, err := sys.GetRequiredEnv("DB_URL")
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Missing DB_URL variable")
	}
	_, err = db.Init(dbUrl, true)
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Error while migrating database: %v", err)
	}
}
