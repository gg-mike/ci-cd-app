package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitProjectDAO(db *gorm.DB) dao.DAO[model.Project, model.ProjectShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name": filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			case "repo": filters["repo LIKE ?"] = "%" + ctx.Query(key) + "%"
			}
		}
	
		return filters, nil
	}

	return dao.DAO[model.Project, model.ProjectShort] {
		DB: db,
		PKCond: func(id uuid.UUID) model.Project { return model.Project { Common: model.Common { ID: id } }},
		Filter: filter,
	}
}
