package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitProjectDAO(db *gorm.DB) dao.DAO[model.Project, model.ProjectCore, model.ProjectCreate, model.ProjectShort] {
	return dao.DAO[model.Project, model.ProjectCore, model.ProjectCreate, model.ProjectShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.ProjectCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.ProjectCreate) error { return nil },
		BeforeUpdate: func(ctx *gin.Context, model model.ProjectCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.ProjectCreate) error { return nil },
		BeforeDelete: func(ctx *gin.Context, model model.Project) error { return errors.New("not implemented") },
		AfterDelete:  func(ctx *gin.Context, model model.Project) error { return nil },
		PKCond:       func(id uuid.UUID) model.Project { return model.Project { Common: model.Common { ID: id } }},
	}
}
