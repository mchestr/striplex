package controllers

import (
	"net/http"
	"plefi/api/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	apicontroller "plefi/api/controllers/api"
	plexcontroller "plefi/api/controllers/plex"
	stripecontroller "plefi/api/controllers/stripe"
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

func (c *AppController) GetRoutes(r *gin.RouterGroup) {
	// Load templates
	r.GET("/health", c.Health)
	r.GET("/logout", c.Logout)

	api := r.Group("/api")
	{
		apiController := apicontroller.NewApiController(api.BasePath(), c.client, c.services)
		apiController.GetRoutes(api)
	}

	plex := r.Group("/plex")
	{
		plexController := plexcontroller.NewPlexController(plex.BasePath(), c.client, c.services)
		plexController.GetRoutes(plex)
	}

	stripe := r.Group("/stripe")
	{
		stripeController := stripecontroller.NewStripeController(stripe.BasePath(), c.client, c.services)
		stripeController.GetRoutes(stripe)
	}
}

// Logout clears the user session by deleting the user_info key
func (p AppController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_info")
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to save session: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully logged out",
	})
}

func (h AppController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Subscriptions displays the subscriptions management page
func (a *AppController) Subscriptions(c *gin.Context) {
	// Check if user is authenticated
	session := sessions.Default(c)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		// If not authenticated, redirect to home page
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Render the subscriptions template
	c.HTML(http.StatusOK, "subscriptions.tmpl", gin.H{})
}
