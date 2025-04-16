package server

import (
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

	// Set up HTML rendering
	router.LoadHTMLGlob("api/views/*")
	router.Static("/static", "./frontend/build/static")
	router.Static("/assets", "./frontend/build/assets")
	router.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")
	router.StaticFile("/manifest.json", "./frontend/build/manifest.json")
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

	// Set Gin mode based on configuration
	gin.SetMode(config.Config.GetString("server.mode"))

	// Initialize session store
	sessionSecret := []byte(config.Config.GetString("server.session_secret"))
	store := cookie.NewStore(sessionSecret)
	router.Use(sessions.Sessions("session_store", store))

	// Initialize controllers
	appController := controllers.NewAppController(client, svcs)
	appController.GetRoutes(&router.RouterGroup)

	return router
}
