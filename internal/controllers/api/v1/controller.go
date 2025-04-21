package v1controller

import (
	"net/http"
	"plefi/internal/config"
	"plefi/internal/middleware"
	"plefi/internal/services"

	"github.com/labstack/echo/v4"
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
func (v *V1) GetRoutes(r *echo.Group) {
	user := r.Group("/user")
	{
		user.GET("/me", middleware.AnonymousHandler(v.GetCurrentUser))
	}

	stripe := r.Group("/stripe")
	{
		stripe.POST("/webhook", v.Webhook)
		// Add new route for subscriptions
		stripe.GET("/subscriptions", middleware.UserHandler(v.GetSubscriptions))
		stripe.POST("/cancel-subscription", middleware.UserHandler(v.CancelUserSubscription))
	}

	plex := r.Group("/plex")
	{
		plex.GET("/check-access", middleware.UserHandler(v.CheckServerAccess))
	}

	// Add new routes for invite code management
	codes := r.Group("/codes")
	codes.Use(middleware.NewAdminMiddleware(config.C.Plex.AdminUserID))
	{
		codes.POST("", v.CreateInviteCode)
		codes.GET("", v.ListInviteCodes)
		codes.GET("/:id", v.GetInviteCode)
		codes.DELETE("/:code", v.DeleteInviteCode)
	}
}
