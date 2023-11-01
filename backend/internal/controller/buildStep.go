package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitBuildStepDAO(db *gorm.DB) dao.DAO[model.BuildStep, model.BuildStepShort] {
	return dao.DAO[model.BuildStep, model.BuildStepShort]{
		DB:     db,
		PKCond: func(id uuid.UUID) model.BuildStep { return model.BuildStep{Common: model.Common{ID: id}} },
		Filter: func(ctx *gin.Context) (map[string]any, error) { return map[string]any{}, nil },
	}
}
