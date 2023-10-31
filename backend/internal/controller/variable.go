package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitVariableDAO(db *gorm.DB) dao.DAO[model.Variable, model.VariableCore, model.VariableCreate, model.VariableShort] {
	return dao.DAO[model.Variable, model.VariableCore, model.VariableCreate, model.VariableShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.VariableCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.VariableCreate) error { return nil },
		BeforeUpdate: func(ctx *gin.Context, model model.VariableCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.VariableCreate) error { return nil },
		BeforeDelete: func(ctx *gin.Context, model model.Variable) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.Variable) error { return nil },
		PKCond:       func(id uuid.UUID) model.Variable { return model.Variable { Common: model.Common { ID: id } }},
	}
}
