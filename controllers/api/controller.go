package apicontroller

import (
	"net/http"
	v1Controller "striplex/controllers/api/v1"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	basePath string
	client   *http.Client
}

func NewApiController(basePath string, client *http.Client) *ApiController {
	return &ApiController{
		basePath: basePath,
		client:   client,
	}
}

func (c *ApiController) GetRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1Controller := v1Controller.NewV1Controller(c.basePath, c.client)
		v1Controller.GetRoutes(v1)
	}
}
