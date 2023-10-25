package main

import (
	_ "github.com/gg-mike/ci-cd-app/backend/docs/server"
	"github.com/gg-mike/ci-cd-app/backend/internal/controllers"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/sys"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
    logger.Basic(zerolog.WarnLevel, "main").Msgf("Missing .env file")
  }

	gin.SetMode(gin.ReleaseMode)

	logDir := sys.GetEnvWithFallback("LOGS", "./logs")
  logger.MultiOutput(logDir)
}

// @title       CI/CD Application - Server API
// @version     1.0.0
// @description Server for CI/CD application developed using Go (Gin, Gorm).
// @basePath    /api
func main() {
	probe := controllers.ReadyProbe{}
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

  rg.GET("/healthz", controllers.Healthz)
  rg.GET("/readyz", probe.Readyz)

  probe.Ready()

  r.Run()
}
