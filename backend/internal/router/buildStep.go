package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary  Get all build steps
// @ID       all-build-steps
// @Tags     build steps
// @Accept   json
// @Produce  json
// @Param    page          query int                 false "Page number"
// @Param    size          query int                 false "Page size"
// @Param    order         query string              false "Order by field"
// @Param    build_step    body  model.BuildStepCore false "Filter"
// @Success  200 {object} []model.BuildStepShort "List of build steps"
// @Failure  400 {object} util.Message           "Error in request"
// @Failure  404 {object} util.Message           "No records found"
// @Failure  500 {object} util.Message           "Database error"
// @Router   /build_steps [get]
func allBuildSteps(dao dao.DAO[model.BuildStep, model.BuildStepCore, model.BuildStepCreate, model.BuildStepShort]) gin.HandlerFunc { return dao.GetMany }

// @Summary  Get the single build step
// @ID       single-build-step
// @Tags     build steps
// @Produce  json
// @Param    id path string true "Build step ID"
// @Success  200 {object} model.BuildStep "Requested build step"
// @Failure  400 {object} util.Message    "Error in params"
// @Failure  404 {object} util.Message    "No record found"
// @Failure  500 {object} util.Message    "Database error"
// @Router   /build_steps/{id} [get]
func getBuildStep(dao dao.DAO[model.BuildStep, model.BuildStepCore, model.BuildStepCreate, model.BuildStepShort]) gin.HandlerFunc { return dao.GetOne }

func InitBuildStepGroup(db *gorm.DB, rg *gin.RouterGroup) {
	dao := controller.InitBuildStepDAO(db)
	
	buildSteps := rg.Group("/build_steps")
	
	buildSteps.GET("", allBuildSteps(dao))
	
	buildSteps.GET("/:id", getBuildStep(dao))
}
