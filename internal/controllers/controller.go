package controllers

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"plefi/internal/config"
	apicontroller "plefi/internal/controllers/api"
	plexcontroller "plefi/internal/controllers/plex"
	stripecontroller "plefi/internal/controllers/stripe"
	"plefi/internal/services"
)

type AppController struct {
	client   *http.Client
	services *services.Services
}

func NewAppController(client *http.Client, services *services.Services) *AppController {
	return &AppController{
		client:   client,
		services: services,
	}
}

func (c *AppController) GetRoutes(r *echo.Echo) {
	r.GET("/health", c.Health)
	r.POST("/logout", c.Logout)

	api := r.Group("/api")
	{
		c := apicontroller.NewApiController("/api", c.client, c.services)
		c.GetRoutes(api)
	}

	plex := r.Group("/plex")
	{
		plexController := plexcontroller.NewPlexController("/plex", c.client, c.services)
		plexController.GetRoutes(plex)
	}

	stripe := r.Group("/stripe")
	{
		stripeController := stripecontroller.NewStripeController("/stripe", c.client, c.services)
		stripeController.GetRoutes(stripe)
	}
}

// Logout clears the user session by deleting the user_info key
func (h AppController) Logout(c echo.Context) error {
	defer c.Redirect(http.StatusFound, "/")
	s, err := session.Get(config.C.Auth.SessionName, c)
	if err != nil {
		return nil
	}
	s.Values["user_info"] = nil
	_ = s.Save(c.Request(), c.Response())
	return nil
}

func (h AppController) Health(c echo.Context) error {
	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
	return echo.NewHTTPError(http.StatusOK, "ok")
}
