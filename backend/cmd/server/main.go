package main

import (
	_ "github.com/gg-mike/ci-cd-app/backend/docs/server"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/router"
	"github.com/gg-mike/ci-cd-app/backend/internal/sys"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Basic(zerolog.WarnLevel, "main").Msgf("Missing .env file")
	}

	gin.SetMode(gin.ReleaseMode)

	logDir := sys.GetEnvWithFallback("LOGS", "./logs")
	logger.MultiOutput("server", logDir)

	dbUrl, err := sys.GetRequiredEnv("DB_URL")
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Missing DB_URL variable")
	}
	DB, err = db.Init(dbUrl, false)
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Error while connecting to database: %v", err)
	}

	vaultAddr, err := sys.GetRequiredEnv("VAULT_ADDR")
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Missing VAULT_ADDR variable")
	}

	vaultRootToken, err := sys.GetRequiredEnv("VAULT_ROOT_TOKEN")
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Missing VAULT_ROOT_TOKEN variable")
	}

	err = vault.Init(vaultAddr, vaultRootToken)
	if err != nil {
		logger.Basic(zerolog.FatalLevel, "main").Msgf("Error while connecting to vault: %v", err)
	}
}

// @title       CI/CD Application - Server API
// @version     1.0.0
// @description Server for CI/CD application developed using Go (Gin, Gorm).
// @basePath    /api
func main() {
	probe := controller.ReadyProbe{}
	probe.Init()
	
	port := sys.GetEnvWithFallback("PORT", "8080")

	r := gin.New()
	r.Use(logger.Gin())
	r.Use(gin.Recovery())

	rg := r.Group("/api")
	rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.
		Basic(zerolog.InfoLevel, "main").
		Msgf("API documentation available at http://localhost:%s/api/docs/index.html", port)

	rg.GET("/healthz", controller.Healthz)
	rg.GET("/readyz", probe.Readyz)

	router.InitBuildGroup(DB, rg)
	router.InitBuildStepGroup(DB, rg)
	router.InitPipelineGroup(DB, rg)
	router.InitProjectGroup(DB, rg)
	router.InitSecretGroup(DB, rg)
	router.InitUserGroup(DB, rg)
	router.InitVariableGroup(DB, rg)
	router.InitWorkerGroup(DB, rg)
	
	probe.Ready()

	r.Run()
}
