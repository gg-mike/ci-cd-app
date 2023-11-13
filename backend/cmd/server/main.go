package main

import (
	"github.com/gg-mike/ci-cd-app/backend/cmd/server/schedulerImpl"
	_ "github.com/gg-mike/ci-cd-app/backend/docs/server"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/router"
	"github.com/gg-mike/ci-cd-app/backend/internal/scheduler"
	"github.com/gg-mike/ci-cd-app/backend/internal/sys"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Warn("main").Msgf("Missing .env file")
	}

	gin.SetMode(gin.ReleaseMode)

	logDir := sys.GetEnvWithFallback("LOGS", "./logs")
	logger.MultiOutput("server", logDir)

	dbUrl, err := sys.GetRequiredEnv("DB_URL")
	if err != nil {
		logger.Fatal("main").Msgf("Missing DB_URL variable")
	}
	if err = db.Init(dbUrl, gorm.Config{Logger: logger.Gorm()}, false); err != nil {
		logger.Fatal("main").Msgf("Error while connecting to database: %v", err)
	}

	vaultAddr, err := sys.GetRequiredEnv("VAULT_ADDR")
	if err != nil {
		logger.Fatal("main").Msgf("Missing VAULT_ADDR variable")
	}

	vaultRootToken, err := sys.GetRequiredEnv("VAULT_ROOT_TOKEN")
	if err != nil {
		logger.Fatal("main").Msgf("Missing VAULT_ROOT_TOKEN variable")
	}

	err = vault.Init(vaultAddr, vaultRootToken)
	if err != nil {
		logger.Fatal("main").Msgf("Error while connecting to vault: %v", err)
	}

	schedulerCtx := schedulerImpl.Context{}
	scheduler.Init(schedulerCtx)
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
		Info("main").
		Msgf("API documentation available at http://localhost:%s/api/docs/index.html", port)

	rg.GET("/healthz", controller.Healthz)
	rg.GET("/readyz", probe.Readyz)

	router.InitBuildGroup(rg)
	router.InitBuildStepGroup(rg)
	router.InitPipelineGroup(rg)
	router.InitProjectGroup(rg)
	router.InitSecretGroup(rg)
	router.InitUserGroup(rg)
	router.InitVariableGroup(rg)
	router.InitWorkerGroup(rg)

	probe.Ready()

	r.Run()
}
