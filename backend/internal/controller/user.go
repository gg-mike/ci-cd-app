package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitUserDAO(db *gorm.DB) dao.DAO[model.User, model.UserCore, model.UserCreate, model.UserShort] {
	return dao.DAO[model.User, model.UserCore, model.UserCreate, model.UserShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.UserCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.UserCreate) error { return errors.New("not implemented") },
		BeforeUpdate: func(ctx *gin.Context, model model.UserCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.UserCreate) error { return errors.New("not implemented") },
		BeforeDelete: func(ctx *gin.Context, model model.User) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.User) error { return errors.New("not implemented") },
		PKCond:       func(id uuid.UUID) model.User { return model.User { Common: model.Common { ID: id } }},
	}
}
