package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitPipelineDAO(db *gorm.DB) dao.DAO[model.Pipeline, model.PipelineShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name":
				filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "branch":
				filters["branch LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "project_id":
				filters["project_id = ?"] = ctx.Query(key)
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Pipeline, model.PipelineShort]{
		DB:     db,
		PKCond: func(id uuid.UUID) model.Pipeline { return model.Pipeline{Common: model.Common{ID: id}} },
		Filter: filter,
	}
}
