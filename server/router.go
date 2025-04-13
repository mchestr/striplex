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
	appController := controllers.NewAppController()

	// Configure routes
	setupAppRoutes(router, appController)
	setupAPIRoutes(router, httpClient)

	return router
}

// setupAppRoutes configures basic application health and info routes.
func setupAppRoutes(router *gin.Engine, appController *controllers.AppController) {
	router.GET("/", appController.Index)
	router.GET("/health", appController.Health)
	router.GET("/whoami", appController.WhoAmI)
}

// setupAPIRoutes configures all API routes with their respective controllers.
func setupAPIRoutes(router *gin.Engine, httpClient *http.Client) {
	// API versioning group
	api := router.Group("/api/v1")

	// Stripe routes
	stripeGroup := api.Group("/stripe")
	stripeController := controllers.NewStripeController(stripeGroup.BasePath())
	stripeGroup.POST("/webhook", stripeController.Webhook)

	// Plex routes
	plexGroup := api.Group("/plex")
	plexController := controllers.NewPlexController(plexGroup.BasePath(), httpClient)
	plexGroup.GET("/auth", plexController.Authenticate)
	plexGroup.GET("/callback", plexController.Callback)
}
