package apicontroller

import (
	"net/http"
	"striplex/config"
	v1Controller "striplex/controllers/api/v1"

	"striplex/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	basePath      string
	client        *http.Client
	wizarrService *services.WizarrService
}

func NewApiController(basePath string, client *http.Client) *ApiController {
	return &ApiController{
		basePath:      basePath,
		client:        client,
		wizarrService: services.NewWizarrService(config.Config.GetString("wizarr.url"), client),
	}
}

func (c *ApiController) GetRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1Controller := v1Controller.NewV1Controller(c.basePath, c.client, c.wizarrService)
		v1Controller.GetRoutes(v1)
	}
}
