package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gin-gonic/gin"
)

// @Summary  Get all pipelines
// @ID       all-pipelines
// @Tags     pipelines
// @Accept   json
// @Produce  json
// @Param    page       query int    false "Page number"
// @Param    size       query int    false "Page size"
// @Param    order      query string false "Order by field"
// @Param    name       query string false "Pipeline name (pattern)"
// @Param    branch     query string false "Pipeline branch (pattern)"
// @Param    project_id query string false "Pipeline project ID (exact)"
// @Success  200 {object} []model.PipelineShort "List of pipelines"
// @Failure  400 {object} dao.Message           "Error in request"
// @Failure  404 {object} dao.Message           "No records found"
// @Failure  500 {object} dao.Message           "Database error"
// @Router   /pipelines [get]
func allPipelines(dao dao.IDAO) gin.HandlerFunc {
	return dao.GetMany
}

// @Summary  Create new pipeline
// @ID       create-pipeline
// @Tags     pipelines
// @Accept   json
// @Produce  json
// @Param    pipeline body model.PipelineInput true "New pipeline entry"
// @Success  201 {object} model.Pipeline "Newly created pipeline"
// @Failure  400 {object} dao.Message    "Error in params"
// @Failure  501 {object} dao.Message    "Endpoint not implemented"
// @Router   /pipelines [post]
func createPipeline(dao dao.IDAO) gin.HandlerFunc {
	return dao.Create
}

// @Summary  Get the single pipeline
// @ID       single-pipeline
// @Tags     pipelines
// @Produce  json
// @Param    id path string true "Pipeline ID"
// @Success  200 {object} model.Pipeline "Requested pipeline"
// @Failure  400 {object} dao.Message    "Error in params"
// @Failure  404 {object} dao.Message    "No record found"
// @Failure  500 {object} dao.Message    "Database error"
// @Router   /pipelines/{id} [get]
func getPipeline(dao dao.IDAO) gin.HandlerFunc { return dao.GetOne }

// @Summary  Update pipeline
// @ID       update-pipeline
// @Tags     pipelines
// @Accept   json
// @Produce  json
// @Param    id       path string               true "Pipeline ID"
// @Param    pipeline body model.PipelineInput true "Updated pipeline entry"
// @Success  200 {object} model.Pipeline "Updated pipeline"
// @Failure  400 {object} dao.Message    "Error in params"
// @Failure  501 {object} dao.Message    "Endpoint not implemented"
// @Router   /pipelines/{id} [put]
func updatePipeline(dao dao.IDAO) gin.HandlerFunc {
	return dao.Update
}

// @Summary  Delete pipeline
// @ID       delete-pipeline
// @Tags     pipelines
// @Produce  json
// @Param    id    path  string  true  "Pipeline ID"
// @Param    force query boolean false "Force deletion"
// @Success  200 {object} dao.Message "Delete message"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  501 {object} dao.Message "Endpoint not implemented"
// @Router   /pipelines/{id} [delete]
func deletePipeline(dao dao.IDAO) gin.HandlerFunc {
	return dao.Delete
}

func InitPipelineGroup(rg *gin.RouterGroup) {
	dao := controller.InitPipelineDAO()

	pipelines := rg.Group("/pipelines")

	pipelines.GET("", allPipelines(dao))
	pipelines.POST("", createPipeline(dao))

	pipelines.GET("/:id", getPipeline(dao))
	pipelines.PUT("/:id", updatePipeline(dao))
	pipelines.DELETE("/:id", deletePipeline(dao))
}
