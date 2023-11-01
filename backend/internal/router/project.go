package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary  Get all projects
// @ID       all-projects
// @Tags     projects
// @Accept   json
// @Produce  json
// @Param    page  query int    false "Page number"
// @Param    size  query int    false "Page size"
// @Param    order query string false "Order by field"
// @Param    name  query string false "Project name (pattern)"
// @Param    repo  query string false "Project repo (pattern)"
// @Success  200 {object} []model.ProjectShort "List of projects"
// @Failure  400 {object} util.Message         "Error in request"
// @Failure  404 {object} util.Message         "No records found"
// @Failure  500 {object} util.Message         "Database error"
// @Router   /projects [get]
func allProjects(dao dao.DAO[model.Project, model.ProjectShort]) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new project
// @ID       create-project
// @Tags     projects
// @Accept   json
// @Produce  json
// @Param    project body model.ProjectCreate true "New project entry"
// @Success  200 {object} model.Project "Newly created project"
// @Failure  400 {object} util.Message  "Error in params"
// @Failure  501 {object} util.Message  "Endpoint not implemented"
// @Router   /projects [post]
func createProject(dao dao.DAO[model.Project, model.ProjectShort]) gin.HandlerFunc { return dao.Create }

// @Summary  Get the single project
// @ID       single-project
// @Tags     projects
// @Produce  json
// @Param    id path string true "Project ID"
// @Success  201 {object} model.Project "Requested project"
// @Failure  400 {object} util.Message  "Error in params"
// @Failure  404 {object} util.Message  "No record found"
// @Failure  500 {object} util.Message  "Database error"
// @Router   /projects/{id} [get]
func getProject(dao dao.DAO[model.Project, model.ProjectShort]) gin.HandlerFunc { return dao.GetOne }

// @Summary  Update project
// @ID       update-project
// @Tags     projects
// @Accept   json
// @Produce  json
// @Param    id      path string              true "Project ID"
// @Param    project body model.ProjectCreate true "Updated project entry"
// @Success  200 {object} model.Project "Updated project"
// @Failure  400 {object} util.Message  "Error in params"
// @Failure  501 {object} util.Message  "Endpoint not implemented"
// @Router   /projects/{id} [put]
func updateProject(dao dao.DAO[model.Project, model.ProjectShort]) gin.HandlerFunc { return dao.Update }

// @Summary  Delete project
// @ID       delete-project
// @Tags     projects
// @Produce  json
// @Param    id    path  string  true  "Project ID"
// @Param    force query boolean false "Force deletion"
// @Success  200 {object} util.Message "Delete message"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /projects/{id} [delete]
func deleteProject(dao dao.DAO[model.Project, model.ProjectShort]) gin.HandlerFunc { return dao.Delete }

func InitProjectGroup(db *gorm.DB, rg *gin.RouterGroup) {
	dao := controller.InitProjectDAO(db)

	projects := rg.Group("/projects")

	projects.GET("", allProjects(dao))
	projects.POST("", createProject(dao))

	projects.GET("/:id", getProject(dao))
	projects.PUT("/:id", updateProject(dao))
	projects.DELETE("/:id", deleteProject(dao))
}
