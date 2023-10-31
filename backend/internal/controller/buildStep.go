package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitBuildStepDAO(db *gorm.DB) dao.DAO[model.BuildStep, model.BuildStepCore, model.BuildStepCreate, model.BuildStepShort] {
	return dao.DAO[model.BuildStep, model.BuildStepCore, model.BuildStepCreate, model.BuildStepShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.BuildStepCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.BuildStepCreate) error { return nil },
		BeforeUpdate: func(ctx *gin.Context, model model.BuildStepCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.BuildStepCreate) error { return nil },
		BeforeDelete: func(ctx *gin.Context, model model.BuildStep) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.BuildStep) error { return nil },
		PKCond:       func(id uuid.UUID) model.BuildStep { return model.BuildStep { Common: model.Common { ID: id } }},
	}
}
