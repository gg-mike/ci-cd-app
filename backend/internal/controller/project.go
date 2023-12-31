package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
)

func InitProjectDAO() dao.IDAO {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name":
				filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "repo":
				filters["repo LIKE ?"] = "%" + ctx.Query(key) + "%"
			}
		}

		return filters, nil
	}

	return dao.DAO[model.Project, model.ProjectShort]{
		Filter: filter,
	}
}
