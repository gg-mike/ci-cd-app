package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func InitSecretDAO() dao.IDAO {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "key":
				filters["key LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "project_id":
				filters["project_id = ?"] = ctx.Query(key)
			case "pipeline_id":
				filters["pipeline_id = ?"] = ctx.Query(key)
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Secret, model.Secret]{
		Filter: filter,
	}
}
