package controller

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitWorkerDAO(db *gorm.DB) dao.DAO[model.Worker, model.WorkerCore, model.WorkerCreate, model.WorkerShort] {
	return dao.DAO[model.Worker, model.WorkerCore, model.WorkerCreate, model.WorkerShort] {
		DB: db,
		BeforeCreate: func(ctx *gin.Context, model model.WorkerCreate) error { return nil },
		AfterCreate:  func(ctx *gin.Context, model model.WorkerCreate) error { return errors.New("not implemented") },
		BeforeUpdate: func(ctx *gin.Context, model model.WorkerCreate) error { return nil },
		AfterUpdate:  func(ctx *gin.Context, model model.WorkerCreate) error { return errors.New("not implemented") },
		BeforeDelete: func(ctx *gin.Context, model model.Worker) error { return nil },
		AfterDelete:  func(ctx *gin.Context, model model.Worker) error { return errors.New("not implemented") },
		PKCond:       func(id uuid.UUID) model.Worker { return model.Worker { Common: model.Common { ID: id } }},
	}
}
