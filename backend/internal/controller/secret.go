package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitSecretDAO(db *gorm.DB) dao.DAO[model.Secret, model.SecretCore, model.SecretCreate, model.SecretShort] {
	return dao.DAO[model.Secret, model.SecretCore, model.SecretCreate, model.SecretShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.SecretCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.SecretCreate) error { return errors.New("not implemented") },
		BeforeUpdate: func(ctx *gin.Context, model model.SecretCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.SecretCreate) error { return errors.New("not implemented") },
		BeforeDelete: func(ctx *gin.Context, model model.Secret) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.Secret) error { return errors.New("not implemented") },
		PKCond:       func(id uuid.UUID) model.Secret { return model.Secret { Common: model.Common { ID: id } }},
	}
}
