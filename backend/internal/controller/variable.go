package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitVariableDAO(db *gorm.DB) dao.DAO[model.Variable, model.VariableShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name": filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			}
		}
	
		return filters, nil
	}

	return dao.DAO[model.Variable, model.VariableShort] {
		DB: db,
		PKCond: func(id uuid.UUID) model.Variable { return model.Variable { Common: model.Common { ID: id } }},
		Filter: filter,
	}
}
