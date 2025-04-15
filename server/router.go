package server

import (
	"net/http"
	"plefi/config"
	"plefi/controllers"
	"plefi/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// NewRouter sets up and configures the application router with all routes.
func NewRouter(client *http.Client) *gin.Engine {
	// Initialize router
	router := gin.Default()

	// Set up HTML rendering
	router.LoadHTMLGlob("views/*")

	// Set Gin mode based on configuration
	gin.SetMode(config.Config.GetString("server.mode"))

	// Initialize session store
	sessionSecret := []byte(config.Config.GetString("server.session_secret"))
	store := cookie.NewStore(sessionSecret)
	router.Use(sessions.Sessions("session_store", store))

	// Initialize controllers
	appController := controllers.NewAppController(client, services.NewServices(client))
	appController.GetRoutes(&router.RouterGroup)

	return router
}
