package stripecontroller

import (
	"fmt"
	"log/slog"
	"net/http"
	"plefi/internal/middleware"
	"plefi/internal/models"
	"plefi/internal/services"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v82"
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
func (s *StripeController) GetRoutes(r *echo.Group) {
	r.GET("/subscribe", middleware.UserHandler(s.CreateCheckoutSession))
	r.GET("/donation", middleware.AnonymousHandler(s.CreateDonationCheckoutSession))
}

// CreateCheckoutSession creates a Stripe checkout session for subscription and redirects the user.
func (h *StripeController) CreateCheckoutSession(c echo.Context, user *models.UserInfo) error {
	customer, err := h.services.Stripe.GetOrCreateCustomer(c.Request().Context(), user)
	if err != nil {
		slog.Error("Failed to get customer for Plex ID", "error", err, "plex_id", user.ID)
		return err
	}
	var anchorDate *time.Time
	if customer != nil {
		sub, err := h.services.Stripe.GetActiveSubscription(c.Request().Context(), user)
		if err != nil {
			slog.Error("Failed to get active subscriptions", "error", err, "plex_id", user.ID)
			return err
		}
		if sub != nil {
			if sub.CancelAtPeriodEnd {
				// If the subscription is set to cancel at the end of the period, we need to set the anchor date
				anchorDateTime := time.Unix(sub.CancelAt, 0)
				anchorDate = &anchorDateTime
			} else {
				return fmt.Errorf("user already has an active subscription")
			}
		}
	}

	// Create or retrieve a customer and checkout session
	sess, err := h.services.Stripe.CreateSubscriptionCheckoutSession(c.Request().Context(), customer, user, anchorDate)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "user", user.Email)
		return err
	}

	// Redirect to Stripe Checkout
	c.Redirect(http.StatusTemporaryRedirect, sess.URL)
	return nil
}

// CreateDonationCheckoutSession creates a Stripe checkout session for donation without requiring authentication
func (h *StripeController) CreateDonationCheckoutSession(c echo.Context, user *models.UserInfo) error {
	var customer *stripe.Customer
	var err error

	// If we have user info, get or create customer
	if user != nil {
		customer, err = h.services.Stripe.GetOrCreateCustomer(c.Request().Context(), user)
		if err != nil {
			slog.Error("Failed to get customer for Plex ID", "error", err, "plex_id", user.ID)
			return err
		}
	}

	// Create donation checkout session
	sess, err := h.services.Stripe.CreateOneTimeCheckoutSession(c.Request().Context(), customer, user)
	if err != nil {
		slog.Error("Failed to create donation checkout session", "error", err)
		return err
	}

	// Redirect to Stripe Checkout
	c.Redirect(http.StatusTemporaryRedirect, sess.URL)
	return nil
}
