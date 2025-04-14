package apicontroller

import (
	"net/http"
	v1Controller "striplex/controllers/api/v1"
	"striplex/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	basePath string
	client   *http.Client
	services *services.Services
}

func NewApiController(basePath string, client *http.Client, services *services.Services) *ApiController {
	return &ApiController{
		basePath: basePath,
		client:   client,
		services: services,
	}
}

func (c *ApiController) GetRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1Controller := v1Controller.NewV1Controller(c.basePath, c.client, c.services)
		v1Controller.GetRoutes(v1)
	}
}
