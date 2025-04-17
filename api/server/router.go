package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"plefi/api/config"
	"plefi/api/controllers"
	"plefi/api/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// NewRouter sets up and configures the application router with all routes.
func NewRouter(svcs *services.Services, client *http.Client) *gin.Engine {
	// Initialize router
	router := gin.Default()

	staticPath := config.C.Server.StaticPath
	slog.Info("static path", "path", staticPath)
	router.Static("/static", fmt.Sprintf("%s/static", staticPath))
	router.Static("/assets", fmt.Sprintf("%s/assets", staticPath))
	router.StaticFile("/favicon.ico", fmt.Sprintf("%s/favicon.ico", staticPath))
	router.StaticFile("/manifest.json", fmt.Sprintf("%s/manifest.json", staticPath))
	router.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", staticPath))
	})

	// Set Gin mode based on configuration
	gin.SetMode(config.C.Server.Mode)

	// Initialize session store
	sessionSecret := []byte(config.C.Auth.SessionSecret)
	store := cookie.NewStore(sessionSecret)
	router.Use(sessions.Sessions("session_store", store))

	// Initialize controllers
	appController := controllers.NewAppController(client, svcs)
	appController.GetRoutes(&router.RouterGroup)

	return router
}
