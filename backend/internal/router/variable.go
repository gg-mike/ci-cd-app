package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary  Get all variables
// @ID       all-variables
// @Tags     variables
// @Accept   json
// @Produce  json
// @Param    page  query int    false "Page number"
// @Param    size  query int    false "Page size"
// @Param    order query string false "Order by field"
// @Param    name  query string false "Variable name (pattern)"
// @Success  200 {object} []model.VariableShort "List of variables"
// @Failure  400 {object} util.Message          "Error in request"
// @Failure  404 {object} util.Message          "No records found"
// @Failure  500 {object} util.Message          "Database error"
// @Router   /variables [get]
func allVariables(dao dao.DAO[model.Variable, model.VariableShort]) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new variable
// @ID       create-variable
// @Tags     variables
// @Accept   json
// @Produce  json
// @Param    variable body model.VariableCreate true "New variable entry"
// @Success  200 {object} model.Variable "Newly created variable"
// @Failure  400 {object} util.Message   "Error in params"
// @Failure  501 {object} util.Message   "Endpoint not implemented"
// @Router   /variables [post]
func createVariable(dao dao.DAO[model.Variable, model.VariableShort]) gin.HandlerFunc { return dao.Create }

// @Summary  Update variable
// @ID       update-variable
// @Tags     variables
// @Accept   json
// @Produce  json
// @Param    id       path string               true "Variable ID"
// @Param    variable body model.VariableCreate true "Updated variable entry"
// @Success  200 {object} model.Variable "Updated variable"
// @Failure  400 {object} util.Message   "Error in params"
// @Failure  501 {object} util.Message   "Endpoint not implemented"
// @Router   /variables/{id} [put]
func updateVariable(dao dao.DAO[model.Variable, model.VariableShort]) gin.HandlerFunc { return dao.Update }

// @Summary  Delete variable
// @ID       delete-variable
// @Tags     variables
// @Produce  json
// @Param    id path string true "Variable ID"
// @Success  200 {object} util.Message "Delete message"
// @Failure  400 {object} util.Message "Error in params"
// @Failure  501 {object} util.Message "Endpoint not implemented"
// @Router   /variables/{id} [delete]
func deleteVariable(dao dao.DAO[model.Variable, model.VariableShort]) gin.HandlerFunc { return dao.Delete }

func InitVariableGroup(db *gorm.DB, rg *gin.RouterGroup) {
	dao := controller.InitVariableDAO(db)
	
	variables := rg.Group("/variables")

	variables.GET( "", allVariables(dao))
	variables.POST("", createVariable(dao))
	
	variables.PUT(   "/:id", updateVariable(dao))
	variables.DELETE("/:id", deleteVariable(dao))
}
