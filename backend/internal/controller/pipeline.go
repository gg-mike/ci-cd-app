package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitPipelineDAO(db *gorm.DB) dao.DAO[model.Pipeline, model.PipelineCore, model.PipelineCreate, model.PipelineShort] {
	return dao.DAO[model.Pipeline, model.PipelineCore, model.PipelineCreate, model.PipelineShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.PipelineCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.PipelineCreate) error { return nil },
		BeforeUpdate: func(ctx *gin.Context, model model.PipelineCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.PipelineCreate) error { return nil },
		BeforeDelete: func(ctx *gin.Context, model model.Pipeline) error { return errors.New("not implemented") },
		AfterDelete:  func(ctx *gin.Context, model model.Pipeline) error { return nil },
		PKCond:       func(id uuid.UUID) model.Pipeline { return model.Pipeline { Common: model.Common { ID: id } }},
	}
}
