package controllers

import (
	"encoding/json"
	"net/http"
	"plefi/config"
	"plefi/model"
	"plefi/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	apicontroller "plefi/controllers/api"
	plexcontroller "plefi/controllers/plex"
	stripecontroller "plefi/controllers/stripe"
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
	r.GET("/whoami", c.WhoAmI)
	r.GET("/logout", c.Logout)
	r.GET("/", c.Index)

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

// WhoAmI returns the authenticated user's information from the session
func (p AppController) WhoAmI(c *gin.Context) {
	userInfo := sessions.Default(c).Get("user_info")
	if userInfo == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Not authenticated",
		})
		return
	}
	var userInfoData model.UserInfo
	err := json.Unmarshal(userInfo.([]byte), &userInfoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to parse user info: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userInfoData,
	})
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

// Index returns a styled HTML landing page for Striplex
func (a *AppController) Index(c *gin.Context) {
	// Get the default price ID from configuration
	priceID := config.Config.GetString("stripe.default_price_id")

	// Check if user is authenticated
	session := sessions.Default(c)
	userInfoData := session.Get("user_info")
	isAuthenticated := userInfoData != nil

	// Prepare template data
	templateData := gin.H{
		"IsAuthenticated": isAuthenticated,
		"PriceID":         priceID,
	}

	// Extract username if authenticated
	if isAuthenticated {
		var userInfo model.UserInfo
		// Handle the case when userInfo is stored as string instead of []byte
		if strData, ok := userInfoData.(string); ok {
			if err := json.Unmarshal([]byte(strData), &userInfo); err == nil {
				templateData["UserInfo"] = userInfo
			}
		}
	}

	// Render the template with data
	c.HTML(http.StatusOK, "index.tmpl", templateData)
}
