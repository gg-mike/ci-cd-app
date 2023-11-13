package router

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/controller"
	"github.com/gg-mike/ci-cd-app/backend/internal/controller/dao"
	"github.com/gin-gonic/gin"
)

// @Summary  Get all secrets
// @ID       all-secrets
// @Tags     secrets
// @Accept   json
// @Produce  json
// @Param    page        query int    false "Page number"
// @Param    size        query int    false "Page size"
// @Param    order       query string false "Order by field"
// @Param    key         query string false "Secret key (pattern)"
// @Param    pipeline_id query string false "Secret project ID (exact)"
// @Param    project_id  query string false "Secret pipeline ID (exact)"
// @Success  200 {object} []model.Secret "List of secrets"
// @Failure  400 {object} dao.Message    "Error in request"
// @Failure  404 {object} dao.Message    "No records found"
// @Failure  500 {object} dao.Message    "Database error"
// @Router   /secrets [get]
func allSecrets(dao dao.IDAO) gin.HandlerFunc { return dao.GetMany }

// @Summary  Create new secret
// @ID       create-secret
// @Tags     secrets
// @Accept   json
// @Produce  json
// @Param    secret body model.SecretInput true "New secret entry"
// @Success  201 {object} model.Secret "Newly created secret"
// @Failure  400 {object} dao.Message  "Error in params"
// @Failure  501 {object} dao.Message  "Endpoint not implemented"
// @Router   /secrets [post]
func createSecret(dao dao.IDAO) gin.HandlerFunc { return dao.Create }

// @Summary  Update secret
// @ID       update-secret
// @Tags     secrets
// @Accept   json
// @Produce  json
// @Param    id     path string            true "Secret ID"
// @Param    secret body model.SecretInput true "Updated secret entry"
// @Success  200 {object} model.Secret "Updated secret"
// @Failure  400 {object} dao.Message  "Error in params"
// @Failure  501 {object} dao.Message  "Endpoint not implemented"
// @Router   /secrets/{id} [put]
func updateSecret(dao dao.IDAO) gin.HandlerFunc { return dao.Update }

// @Summary  Delete secret
// @ID       delete-secret
// @Tags     secrets
// @Produce  json
// @Param    id path string true "Secret ID"
// @Success  200 {object} dao.Message "Delete message"
// @Failure  400 {object} dao.Message "Error in params"
// @Failure  501 {object} dao.Message "Endpoint not implemented"
// @Router   /secrets/{id} [delete]
func deleteSecret(dao dao.IDAO) gin.HandlerFunc { return dao.Delete }

func InitSecretGroup(rg *gin.RouterGroup) {
	dao := controller.InitSecretDAO()

	secrets := rg.Group("/secrets")

	secrets.GET("", allSecrets(dao))
	secrets.POST("", createSecret(dao))

	secrets.PUT("/:id", updateSecret(dao))
	secrets.DELETE("/:id", deleteSecret(dao))
}
