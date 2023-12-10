package queue

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gin-gonic/gin"
)

func DAO() dao.IDAO {
	return dao.DAO[QueueElem, QueueElem]{
		Filter: func(ctx *gin.Context) (map[string]any, error) { return map[string]any{}, nil },
	}
}
