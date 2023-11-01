package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary  Get all workers
// @ID       all-workers
// @Tags     workers
// @Accept   json
// @Produce  json
// @Param    page   query int    false "Page number"
// @Param    size   query int    false "Page size"
// @Param    order  query string false "Order by field"
// @Param    name   query string false "Worker name (pattern)"
// @Param    system query string false "Worker system (pattern)"
// @Param    status query []int  false "Worker status (possible values)"
// @Param    type   query []int  false "Worker type (possible values)"
// @Success  200 {object} []model.WorkerShort "List of workers"
// @Failure  400 {object} util.Message        "Error in request"
// @Failure  404 {object} util.Message        "No records found"
// @Failure  500 {object} util.Message        "Database error"
// @Router   /workers [get]
func allWorkers(dao dao.DAO[model.Worker, model.WorkerShort]) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new worker
// @ID       create-worker
// @Tags     workers
// @Accept   json
// @Produce  json
// @Param    worker body model.WorkerCreate true "New worker entry"
// @Success  200 {object} model.Worker "Newly created worker"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /workers [post]
func createWorker(dao dao.DAO[model.Worker, model.WorkerShort]) gin.HandlerFunc { return dao.Create }

// @Summary  Get the single worker
// @ID       single-worker
// @Tags     workers
// @Produce  json
// @Param    id path string true "Worker ID"
// @Success  200 {object} model.Worker "Requested worker"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  404 {object} util.Message "No record found"
// @Failure  500 {object} util.Message "Database error"
// @Router   /workers/{id} [get]
func getWorker(dao dao.DAO[model.Worker, model.WorkerShort]) gin.HandlerFunc { return dao.GetOne }

// @Summary  Update worker
// @ID       update-worker
// @Tags     workers
// @Accept   json
// @Produce  json
// @Param    id     path string             true "Worker ID"
// @Param    worker body model.WorkerCreate true "Updated worker entry"
// @Success  200 {object} model.Worker "Updated worker"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /workers/{id} [put]
func updateWorker(dao dao.DAO[model.Worker, model.WorkerShort]) gin.HandlerFunc { return dao.Update }

// @Summary  Delete worker
// @ID       delete-worker
// @Tags     workers
// @Produce  json
// @Param    id path string true "Worker ID"
// @Success  200 {object} util.Message "Delete message"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /workers/{id} [delete]
func deleteWorker(dao dao.DAO[model.Worker, model.WorkerShort]) gin.HandlerFunc { return dao.Delete }

func InitWorkerGroup(db *gorm.DB, rg *gin.RouterGroup) {
	dao := controller.InitWorkerDAO(db)

	workers := rg.Group("/workers")

	workers.GET("", allWorkers(dao))
	workers.POST("", createWorker(dao))

	workers.GET("/:id", getWorker(dao))
	workers.PUT("/:id", updateWorker(dao))
	workers.DELETE("/:id", deleteWorker(dao))
}
