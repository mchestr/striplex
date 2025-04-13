package server

import (
	"net/http"
	"striplex/config"
	"striplex/controllers"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// NewRouter sets up and configures the application router with all routes.
func NewRouter() *gin.Engine {
	// Set Gin mode based on configuration
	gin.SetMode(config.Config.GetString("server.mode"))

	// Create HTTP client with reasonable timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Initialize router and session store
	router := gin.Default()
	sessionSecret := []byte(config.Config.GetString("server.session_secret"))
	store := cookie.NewStore(sessionSecret)
	router.Use(sessions.Sessions("session_store", store))

	// Initialize controllers
	appController := controllers.NewAppController(httpClient)
	appController.GetRoutes(&router.RouterGroup)

	return router
}
