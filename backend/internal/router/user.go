package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gin-gonic/gin"
)

// @Summary  Get all users
// @ID       all-users
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    page     query int    false "Page number"
// @Param    size     query int    false "Page size"
// @Param    order    query string false "Order by field"
// @Param    username query string false "Username (pattern)"
// @Success  200 {object} []model.User "List of users"
// @Failure  400 {object} dao.Message  "Error in request"
// @Failure  404 {object} dao.Message  "No records found"
// @Failure  500 {object} dao.Message  "Database error"
// @Router   /users [get]
func allUsers(dao dao.IDAO) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new user
// @ID       create-user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    user body model.UserInput true "New user entry"
// @Success  200 {object} model.User  "Newly created user"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  501 {object} dao.Message "Endpoint not implemented"
// @Router   /users [post]
func createUser(dao dao.IDAO) gin.HandlerFunc { return dao.Create }

// @Summary  Get the single user
// @ID       single-user
// @Tags     users
// @Produce  json
// @Param    id path string true "User ID"
// @Success  200 {object} model.User  "Requested user"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  404 {object} dao.Message "No record found"
// @Failure  500 {object} dao.Message "Database error"
// @Router   /users/{id} [get]
func getUser(dao dao.IDAO) gin.HandlerFunc { return dao.GetOne }

// @Summary  Update user
// @ID       update-user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    id   path string          true "User ID"
// @Param    user body model.UserInput true "Updated user entry"
// @Success  200 {object} model.User  "Updated user"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  501 {object} dao.Message "Endpoint not implemented"
// @Router   /users/{id} [put]
func updateUser(dao dao.IDAO) gin.HandlerFunc { return dao.Update }

// @Summary  Delete user
// @ID       delete-user
// @Tags     users
// @Produce  json
// @Param    id path string true "User ID"
// @Success  200 {object} dao.Message "Delete message"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  501 {object} dao.Message "Endpoint not implemented"
// @Router   /users/{id} [delete]
func deleteUser(dao dao.IDAO) gin.HandlerFunc { return dao.Delete }

func InitUserGroup(rg *gin.RouterGroup) {
	dao := controller.InitUserDAO()

	users := rg.Group("/users")

	users.GET("", allUsers(dao))
	users.POST("", createUser(dao))

	users.GET("/:id", getUser(dao))
	users.PUT("/:id", updateUser(dao))
	users.DELETE("/:id", deleteUser(dao))
}
