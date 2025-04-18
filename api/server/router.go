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
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.C.Auth.SessionSecret))))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  config.C.Server.StaticPath,
		HTML5: true,
	}))

	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(true),                          // e.g. ipv4 start with 127.
		echo.TrustIPRange(config.C.Server.TrustedProxies), // use parsed CIDRs
	)

	// Initialize controllers
	appController := controllers.NewAppController(client, svcs)
	appController.GetRoutes(e)
	return e
}
