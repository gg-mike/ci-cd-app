package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitBuildDAO(db *gorm.DB) dao.DAO[model.Build, model.BuildCore, model.BuildCreate, model.BuildShort] {
	return dao.DAO[model.Build, model.BuildCore, model.BuildCreate, model.BuildShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.BuildCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.BuildCreate) error { return errors.New("not implemented") },
		BeforeUpdate: func(ctx *gin.Context, model model.BuildCreate) error { return errors.New("not implemented") },
		AfterUpdate:  func(ctx *gin.Context, model model.BuildCreate) error { return errors.New("not implemented") },
		BeforeDelete: func(ctx *gin.Context, model model.Build) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.Build) error { return nil },
		PKCond:       func(id uuid.UUID) model.Build { return model.Build { Common: model.Common { ID: id } }},
	}
}
