package main

import (
	_ "github.com/gg-mike/ci-cd-app/backend/docs/server"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/sys"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("main").Msgf("Missing .env file")
	}

	logDir := sys.GetEnvWithFallback("LOGS", "./logs")
	logger.MultiOutput("dbMigrator", logDir)

	dbUrl, err := sys.GetRequiredEnv("DB_URL")
	if err != nil {
		logger.Fatal("main").Msgf("Missing DB_URL variable")
	}

	if err = db.Init(dbUrl, gorm.Config{Logger: logger.Gorm()}, true); err != nil {
		logger.Fatal("main").Msgf("Error while migrating database: %v", err)
	}
}
