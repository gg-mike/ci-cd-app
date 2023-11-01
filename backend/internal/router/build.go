package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary  Get all builds
// @ID       all-builds
// @Tags     builds
// @Accept   json
// @Produce  json
// @Param    page        query int    false "Page number"
// @Param    size        query int    false "Page size"
// @Param    order       query string false "Order by field"
// @Param    status      query []int  false "Build status (possible values)"
// @Param    worker_id   query string false "Build worker ID (exact)"
// @Param    pipeline_id query string false "Build pipeline ID (exact)"
// @Success  200 {object} []model.BuildShort "List of builds"
// @Failure  400 {object} util.Message       "Error in request"
// @Failure  404 {object} util.Message       "No records found"
// @Failure  500 {object} util.Message       "Database error"
// @Router   /builds [get]
func allBuilds(dao dao.DAO[model.Build, model.BuildShort]) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new build
// @ID       create-build
// @Tags     builds
// @Accept   json
// @Produce  json
// @Param    build body model.BuildCreate true "New build entry"
// @Success  201 {object} model.Build  "Newly created build"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /builds [post]
func createBuild(dao dao.DAO[model.Build, model.BuildShort]) gin.HandlerFunc { return dao.Create }

// @Summary  Get the single build
// @ID       single-build
// @Tags     builds
// @Produce  json
// @Param    id path string true "Build ID"
// @Success  200 {object} model.Build  "Requested build"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  404 {object} util.Message "No record found"
// @Failure  500 {object} util.Message "Database error"
// @Router   /builds/{id} [get]
func getBuild(dao dao.DAO[model.Build, model.BuildShort]) gin.HandlerFunc { return dao.GetOne }

// @Summary  Update build
// @ID       update-build
// @Tags     builds
// @Accept   json
// @Produce  json
// @Param    id    path string            true "Build ID"
// @Param    build body model.BuildCreate true "Updated build entry"
// @Success  200 {object} model.Build  "Updated build"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /builds/{id} [put]
func updateBuild(dao dao.DAO[model.Build, model.BuildShort]) gin.HandlerFunc { return dao.Update }

func InitBuildGroup(db *gorm.DB, rg *gin.RouterGroup) {
	dao := controller.InitBuildDAO(db)

	builds := rg.Group("/builds")

	builds.GET("", allBuilds(dao))
	builds.POST("", createBuild(dao))

	builds.GET("/:id", getBuild(dao))
	builds.PUT("/:id", updateBuild(dao))
}
