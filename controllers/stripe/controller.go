package stripecontroller

import (
	"fmt"
	"log/slog"
	"net/http"
	"plefi/config"
	"plefi/model"
	"plefi/services"

	"github.com/gin-gonic/gin"
)

// StripeController handles Stripe payment and subscription related operations
type StripeController struct {
	basePath string
	client   *http.Client
	services *services.Services
}

// NewStripeController creates a new StripeController instance
func NewStripeController(basePath string, client *http.Client, services *services.Services) *StripeController {
	return &StripeController{
		basePath: basePath,
		client:   client,
		services: services,
	}
}

// GetRoutes registers all routes for the Stripe controller
func (s *StripeController) GetRoutes(r *gin.RouterGroup) {
	r.GET("/checkout", s.CreateCheckoutSession)
	r.GET("/success", s.SuccessSubscription)
	r.GET("/cancel", s.CancelSubscription)
}

// CreateCheckoutSession creates a Stripe checkout session for subscription and redirects the user.
func (s *StripeController) CreateCheckoutSession(ctx *gin.Context) {
	// Check for Plex authentication in session
	userInfo, err := model.GetUserInfo(ctx)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}
	if userInfo == nil {
		// Redirect to Plex authentication route
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/plex/auth?next=%s/checkout",
			s.basePath))
		return
	}

	customer, err := s.services.Stripe.GetOrCreateCustomer(ctx, userInfo)
	if err != nil {
		slog.Error("Failed to get customer for Plex ID", "error", err, "plex_id", userInfo.ID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer"})
		return
	}

	// Create or retrieve a customer and checkout session
	sess, err := s.services.Stripe.CreateSubscriptionCheckoutSession(ctx, customer, userInfo)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "user", userInfo.Email)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
		return
	}

	// Redirect to Stripe Checkout
	ctx.Redirect(http.StatusTemporaryRedirect, sess.URL)
}

// SuccessSubscription handles successful Stripe checkout
func (s *StripeController) SuccessSubscription(ctx *gin.Context) {
	userInfo, err := model.GetUserInfo(ctx)
	if err != nil || userInfo == nil {
		slog.Error("Failed to parse user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session data"})
		return
	}

	// Render the template with data
	ctx.HTML(http.StatusOK, "stripe_success.tmpl", gin.H{
		"UserInfo": userInfo,
	})
}

// CancelSubscription handles cancelled Stripe checkout
func (s *StripeController) CancelSubscription(ctx *gin.Context) {
	priceID := ctx.Query("price_id")
	if priceID == "" {
		priceID = config.Config.GetString("stripe.default_price_id")
	}

	// Render the cancel template with data
	ctx.HTML(http.StatusOK, "stripe_cancel.tmpl", gin.H{
		"PriceID": priceID,
	})
}
