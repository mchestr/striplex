package server

import (
	"net/http"
	"plefi/api/config"
	"plefi/api/controllers"
	"plefi/api/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter sets up and configures the application router with all routes.
func NewRouter(svcs *services.Services, client *http.Client) *echo.Echo {
	// Initialize router
	r := echo.New()
	r.Use(middleware.Recover())
	r.Use(middleware.Logger())
	r.Use(session.Middleware(sessions.NewCookieStore([]byte(config.C.Auth.SessionSecret))))
	r.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  config.C.Server.StaticPath,
		HTML5: true,
	}))

	// Initialize controllers
	appController := controllers.NewAppController(client, svcs)
	appController.GetRoutes(r)
	return r
}
