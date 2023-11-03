package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func InitWorkerDAO() dao.DAO[model.Worker, model.WorkerShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name":
				filters["name   LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "system":
				filters["system LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "status":
				filters["status IN ?"] = ctx.QueryArray(key)
			case "type":
				filters["type   IN ?"] = ctx.QueryArray(key)
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Worker, model.WorkerShort]{
		Filter: filter,
	}
}
