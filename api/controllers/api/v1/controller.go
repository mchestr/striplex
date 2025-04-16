package v1controller

import (
	"net/http"
	"plefi/api/services"

	"github.com/gin-gonic/gin"
)

type V1 struct {
	basePath string
	client   *http.Client
	services *services.Services
}

func NewV1Controller(basePath string, client *http.Client, services *services.Services) *V1 {
	return &V1{
		basePath: basePath,
		client:   client,
		services: services,
	}
}
func (v *V1) GetRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.GET("/me", v.GetCurrentUser)
	}

	stripe := r.Group("/stripe")
	{
		stripe.POST("/webhook", v.Webhook)
		// Add new route for subscriptions
		stripe.GET("/subscriptions", v.GetSubscriptions)
		stripe.POST("/cancel-subscription", v.CancelUserSubscription)
	}

	plex := r.Group("/plex")
	{
		plex.GET("/check-access", v.CheckServerAccess)
	}
}
