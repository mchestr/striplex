package v1controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type V1 struct {
	basePath string
	client   *http.Client
}

func NewV1Controller(basePath string, client *http.Client) *V1 {
	return &V1{
		basePath: basePath,
		client:   client,
	}
}
func (v *V1) GetRoutes(r *gin.RouterGroup) {
	stripe := r.Group("/stripe")
	{
		stripe.POST("/webhook", v.Webhook)
	}
}
