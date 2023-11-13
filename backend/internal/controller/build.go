package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitBuildDAO() dao.IDAO {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "status":
				filters["status in ?"] = ctx.QueryArray(key)
			case "worker_id":
				filters["worker_id = ?"] = ctx.Query(key)
			case "pipeline_id":
				filters["pipeline_id = ?"] = ctx.Query(key)
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Build, model.BuildShort]{
		Preload: func(tx *gorm.DB) *gorm.DB { return tx.Preload("Steps.Logs").Preload(clause.Associations) },
		Filter:  filter,
	}
}
