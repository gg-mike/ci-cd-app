package controller

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitSecretDAO(db *gorm.DB) dao.DAO[model.Secret, model.SecretShort] {
	filter := func(ctx *gin.Context) (map[string]any, error) {
		filters := map[string]any{}
		for key := range ctx.Request.URL.Query() {
			switch key {
			case "name": filters["name LIKE ?"] = "%" + ctx.Query(key) + "%"
			}
		}
	
		return filters, nil
	}

	return dao.DAO[model.Secret, model.SecretShort] {
		DB: db,
		PKCond: func(id uuid.UUID) model.Secret { return model.Secret { Common: model.Common { ID: id } }},
		Filter: filter,
	}
}
