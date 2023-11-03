package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func InitVariableDAO() dao.DAO[model.Variable, model.Variable] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name":
				filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Variable, model.Variable]{
		Filter: filter,
	}
}
