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
	adminMiddleware := middleware.NewAdminMiddleware(config.C.Plex.AdminUserID)

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
		admin := plex.Group("/users", adminMiddleware)
		{
			admin.GET("", v.GetPlexUsers)
			admin.GET("/:id", v.GetPlexUser)
			admin.GET("/:id/invites", v.GetPlexUserInvites)
		}
		plex.GET("/check-access", middleware.UserHandler(v.CheckServerAccess))
	}
	// Add new routes for invite code management
	codes := r.Group("/codes")
	{
		admin := codes.Group("", adminMiddleware)
		{
			admin.POST("", v.CreateInviteCode)
			admin.GET("", v.ListInviteCodes)
			admin.GET("/:id", v.GetInviteCode)
			admin.DELETE("/:code", v.DeleteInviteCode)
		}
		codes.POST("/claim", middleware.UserHandler(v.ClaimInviteCode))
	}
}
