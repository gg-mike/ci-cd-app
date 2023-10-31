package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitWorkerDAO(db *gorm.DB) dao.DAO[model.Worker, model.WorkerShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name":   filters["name   LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "system": filters["system LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "status": filters["status IN ?"]   = ctx.QueryArray(key)
			case "type":   filters["type   IN ?"]   = ctx.QueryArray(key)
			}
		}
	
		return filters, nil
	}

	return dao.DAO[model.Worker, model.WorkerShort] {
		DB: db,
		PKCond: func(id uuid.UUID) model.Worker { return model.Worker { Common: model.Common { ID: id } }},
		Filter: filter,
	}
}
