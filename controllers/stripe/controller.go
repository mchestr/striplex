package stripecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"plefi/config"
	"plefi/model"
	"plefi/services"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	stripeSession "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
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
	r.POST("/cancel-subscription", s.CancelUserSubscription)
}

// CreateCheckoutSession creates a Stripe checkout session for subscription and redirects the user.
func (s *StripeController) CreateCheckoutSession(ctx *gin.Context) {
	// Check for Plex authentication in session
	session := sessions.Default(ctx)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		// If no user info is found, redirect to Plex auth flow
		// Store original request info to return after auth
		redirectURL := ctx.Request.URL.String()
		session.Set("checkout_redirect", redirectURL)
		if err := session.Save(); err != nil {
			slog.Error("Failed to save session", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		price_id := ctx.Query("price_id")
		if price_id == "" {
			price_id = config.Config.GetString("stripe.default_price_id")
		}

		// Redirect to Plex authentication route
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/plex/auth?next=%s/checkout?price_id=%s",
			s.basePath, price_id))
		return
	}

	// Parse the Plex user info
	userInfoData, err := parseUserInfo(userInfo)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}

	// Get the price ID from the request (could be query param or from body)
	priceID := ctx.Query("price_id")
	if priceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing price_id parameter"})
		return
	}

	// Create or retrieve a customer and checkout session
	sess, err := s.createCheckoutSession(ctx, userInfoData, priceID)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "user", userInfoData.Email)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
		return
	}

	// Redirect to Stripe Checkout
	ctx.Redirect(http.StatusSeeOther, sess.URL)
}

// SuccessSubscription handles successful Stripe checkout
func (s *StripeController) SuccessSubscription(ctx *gin.Context) {
	// Get user info from session if available
	session := sessions.Default(ctx)
	userInfo := session.Get("user_info")

	// Parse user info if available
	var plexUser *model.UserInfo
	var err error
	if userInfo != nil {
		plexUser, err = parseUserInfo(userInfo)
		if err != nil {
			slog.Warn("failed to parse user info in success page", "error", err)
		}
	}

	// Render the template with data
	ctx.HTML(http.StatusOK, "stripe_success.tmpl", gin.H{
		"Username": plexUser.Username,
		"Email":    plexUser.Email,
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

// CancelUserSubscription cancels a user's active subscription
func (s *StripeController) CancelUserSubscription(ctx *gin.Context) {
	// Check for authentication
	session := sessions.Default(ctx)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		// Redirect to Plex authentication route
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/plex/auth?next=%s/cancel",
			s.basePath))
		return
	}

	// Parse user info
	plexUser, err := parseUserInfo(userInfo)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Invalid session data",
		})
		return
	}

	// Find Stripe customer by Plex user ID
	customerList := &stripe.CustomerListParams{}
	customerList.Filters.AddFilter("metadata[plex_user_id]", "", strconv.Itoa(plexUser.ID))

	it := customer.List(customerList)

	// Track if we found and processed any customers
	customersFound := false
	subscriptionsCanceled := 0

	// Process each customer found
	for it.Next() {
		customersFound = true
		cus := it.Customer()

		// Find active subscriptions for this customer
		subList := &stripe.SubscriptionListParams{}
		subList.Filters.AddFilter("customer", "", cus.ID)
		subList.Filters.AddFilter("status", "", "active")

		subIt := subscription.List(subList)

		// Process each subscription for this customer
		for subIt.Next() {
			sub := subIt.Subscription()

			// Cancel subscription
			cancelParams := &stripe.SubscriptionParams{
				CancelAtPeriodEnd: stripe.Bool(true),
			}

			_, err = subscription.Update(sub.ID, cancelParams)
			if err != nil {
				slog.Error("Failed to cancel subscription",
					"error", err,
					"subscription_id", sub.ID,
					"customer_id", cus.ID)
				continue
			}

			subscriptionsCanceled++
			slog.Info("Canceled subscription",
				"subscription_id", sub.ID,
				"customer_id", cus.ID,
				"plex_user_id", plexUser.ID)
		}
	}

	// Check if we processed any customers
	if !customersFound {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "No customers found for this user",
		})
		return
	}

	// Check if we canceled any subscriptions
	if subscriptionsCanceled == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "No active subscriptions found for this user",
		})
		return
	}

	// Return success
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Successfully canceled %d subscription(s)", subscriptionsCanceled),
	})
}

// Helper functions

// createCheckoutSession creates a Stripe customer and checkout session
func (s *StripeController) createCheckoutSession(ctx context.Context, userData *model.UserInfo, priceID string) (*stripe.CheckoutSession, error) {
	// Set success and cancel URLs
	hostname := config.Config.GetString("server.hostname")
	successURL := fmt.Sprintf("https://%s%s/success", hostname, s.basePath)
	cancelURL := fmt.Sprintf("https://%s%s/cancel?price_id=%s", hostname, s.basePath, priceID)

	// Create a Stripe customer first with Plex user metadata
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(userData.Email),
		Name:  stripe.String(userData.Username),
		Metadata: map[string]string{
			"plex_user_id":  strconv.Itoa(userData.ID),
			"plex_username": userData.Username,
			"plex_email":    userData.Email,
		},
	}
	customerParams.Context = ctx

	customer, err := customer.New(customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	slog.Info("Created Stripe customer",
		"customer_id", customer.ID,
		"email", customer.Email,
		"plex_id", userData.ID)

	// Create checkout session parameters
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Customer:   stripe.String(customer.ID),
	}
	params.Context = ctx

	// Create the checkout session
	return stripeSession.New(params)
}

// parseUserInfo parses user info from the session
func parseUserInfo(userInfo interface{}) (*model.UserInfo, error) {
	var userInfoData model.UserInfo

	if byteData, ok := userInfo.(string); ok {
		if err := json.Unmarshal([]byte(byteData), &userInfoData); err != nil {
			return nil, fmt.Errorf("invalid user info JSON: %w", err)
		}
		return &userInfoData, nil
	}

	return nil, fmt.Errorf("user info is not in expected string format")
}
