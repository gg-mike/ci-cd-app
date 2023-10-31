package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitBuildDAO(db *gorm.DB) dao.DAO[model.Build, model.BuildShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "status":      filters["status in ?"]     = ctx.QueryArray(key)
			case "worker_id":   filters["worker_id = ?"]   = ctx.Query(key)
			case "pipeline_id": filters["pipeline_id = ?"] = ctx.Query(key)	
			}
		}
	
		return filters, nil
	}

	return dao.DAO[model.Build, model.BuildShort] {
		DB: db,
		PKCond: func(id uuid.UUID) model.Build { return model.Build { Common: model.Common { ID: id } }},
		Filter: filter,
	}
}
