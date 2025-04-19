package stripecontroller

import (
	"fmt"
	"log/slog"
	"net/http"
	"plefi/internal/models"
	"plefi/internal/services"
	"plefi/internal/utils"
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
	r.GET("/subscribe", s.CreateCheckoutSession)
	r.GET("/donation", s.CreateDonationCheckoutSession)
}

// CreateCheckoutSession creates a Stripe checkout session for subscription and redirects the user.
func (h *StripeController) CreateCheckoutSession(c echo.Context) error {
	// Check for Plex authentication in session
	userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		return err
	}
	if userInfo == nil {
		// Redirect to Plex authentication route
		c.Redirect(http.StatusFound, fmt.Sprintf("/plex/auth?next=%s/subscribe",
			h.basePath))
		return nil
	}
	userInfoData, ok := userInfo.(*models.UserInfo)
	if !ok {
		slog.Error("Failed to cast user info to UserInfo type")
		return fmt.Errorf("failed to cast user info to UserInfo type")
	}

	customer, err := h.services.Stripe.GetOrCreateCustomer(c.Request().Context(), userInfoData)
	if err != nil {
		slog.Error("Failed to get customer for Plex ID", "error", err, "plex_id", userInfoData.ID)
		return err
	}
	var anchorDate *time.Time
	if customer != nil {
		sub, err := h.services.Stripe.GetActiveSubscription(c.Request().Context(), userInfoData)
		if err != nil {
			slog.Error("Failed to get active subscriptions", "error", err, "plex_id", userInfoData.ID)
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
	sess, err := h.services.Stripe.CreateSubscriptionCheckoutSession(c.Request().Context(), customer, userInfoData, anchorDate)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "user", userInfoData.Email)
		return err
	}

	// Redirect to Stripe Checkout
	c.Redirect(http.StatusTemporaryRedirect, sess.URL)
	return nil
}

// CreateDonationCheckoutSession creates a Stripe checkout session for donation without requiring authentication
func (h *StripeController) CreateDonationCheckoutSession(c echo.Context) error {
	var customer *stripe.Customer
	var err error

	// Try to get user info if available, but don't require it
	userInfo, _ := utils.GetSessionData(c, utils.UserInfoState)
	userInfoData, _ := userInfo.(*models.UserInfo)

	// If we have user info, get or create customer
	if userInfoData != nil {
		customer, err = h.services.Stripe.GetOrCreateCustomer(c.Request().Context(), userInfoData)
		if err != nil {
			slog.Error("Failed to get customer for Plex ID", "error", err, "plex_id", userInfoData.ID)
			return err
		}
	}

	// Create donation checkout session
	sess, err := h.services.Stripe.CreateOneTimeCheckoutSession(c.Request().Context(), customer, userInfoData)
	if err != nil {
		slog.Error("Failed to create donation checkout session", "error", err)
		return err
	}

	// Redirect to Stripe Checkout
	c.Redirect(http.StatusTemporaryRedirect, sess.URL)
	return nil
}
