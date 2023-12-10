package queue

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gin-gonic/gin"
)

// @Summary  Get all queue elems
// @ID       all-queue-elems
// @Tags     queue-elem
// @Accept   json
// @Produce  json
// @Param    page        query int    false "Page number"
// @Param    size        query int    false "Page size"
// @Param    order       query string false "Order by field"
// @Success  200 {object} []QueueElem "List of queue elems"
// @Failure  400 {object} dao.Message "Error in request"
// @Failure  404 {object} dao.Message "No records found"
// @Failure  500 {object} dao.Message "Database error"
// @Router   /queue-elems [get]
func allQueueElems(dao dao.IDAO) gin.HandlerFunc { return dao.GetMany }

func Router(rg *gin.RouterGroup) {
	dao := DAO()

	builds := rg.Group("/queue-elems")

	builds.GET("", allQueueElems(dao))
}
