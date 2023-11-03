package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func InitBuildStepDAO() dao.DAO[model.BuildStep, model.BuildStepShort] {
	return dao.DAO[model.BuildStep, model.BuildStepShort]{
		Filter: func(ctx *gin.Context) (map[string]any, error) { return map[string]any{}, nil },
	}
}
