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
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	stripeSession "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
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
func (s *StripeController) CreateCheckoutSession(c *gin.Context) {
	// Create a context with timeout for external API calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Check for Plex authentication in session
	session := sessions.Default(c)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		// If no user info is found, redirect to Plex auth flow
		// Store original request info to return after auth
		redirectURL := c.Request.URL.String()
		session.Set("checkout_redirect", redirectURL)
		if err := session.Save(); err != nil {
			slog.Error("Failed to save session", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		// Redirect to Plex authentication route
		c.Redirect(http.StatusFound, fmt.Sprintf("/plex/auth?next=%s/checkout?price_id=%s",
			s.basePath, c.Query("price_id")))
		return
	}

	// Parse the Plex user info
	userInfoData, err := parseUserInfo(userInfo)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data"})
		return
	}

	// Get the price ID from the request (could be query param or from body)
	priceID := c.Query("price_id")
	if priceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing price_id parameter"})
		return
	}

	// Create or retrieve a customer and checkout session
	sess, err := s.createCheckoutSession(ctx, userInfoData, priceID)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err, "user", userInfoData.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
		return
	}

	// Redirect to Stripe Checkout
	c.Redirect(http.StatusSeeOther, sess.URL)
}

// SuccessSubscription handles successful Stripe checkout
func (s *StripeController) SuccessSubscription(c *gin.Context) {
	// Get user info from session if available
	session := sessions.Default(c)
	userInfo := session.Get("user_info")

	// Parse user info if available
	var plexUser model.UserInfo
	if userInfo != nil {
		parsedUser, err := parseUserInfo(userInfo)
		if err != nil {
			slog.Warn("failed to parse user info in success page", "error", err)
		} else {
			plexUser = parsedUser
		}
	}

	// Generate dynamic content based on available user info
	username := "Thank you!"
	if plexUser.Username != "" {
		username = "Thank you, " + plexUser.Username + "!"
	}

	emailDisplay := "your email address"
	if plexUser.Email != "" {
		emailDisplay = "<strong>" + plexUser.Email + "</strong>"
	}

	// Display a success page
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, generateSuccessHTML(username, emailDisplay))
}

// CancelSubscription handles cancelled Stripe checkout
func (s *StripeController) CancelSubscription(c *gin.Context) {
	priceID := c.Query("price_id")

	// Display a cancellation page
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, generateCancelHTML(priceID))
}

// Helper functions

// createCheckoutSession creates a Stripe customer and checkout session
func (s *StripeController) createCheckoutSession(ctx context.Context, userData model.UserInfo, priceID string) (*stripe.CheckoutSession, error) {
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
func parseUserInfo(userInfo interface{}) (model.UserInfo, error) {
	var userInfoData model.UserInfo

	if byteData, ok := userInfo.(string); ok {
		if err := json.Unmarshal([]byte(byteData), &userInfoData); err != nil {
			return model.UserInfo{}, fmt.Errorf("invalid user info JSON: %w", err)
		}
		return userInfoData, nil
	}

	return model.UserInfo{}, fmt.Errorf("user info is not in expected string format")
}

// generateSuccessHTML generates the HTML for the success page
func generateSuccessHTML(username, emailDisplay string) string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Subscription Success</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #1e1e2e;
            color: #cdd6f4;
            margin: 0;
            padding: 0;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
        }
        .container {
            text-align: center;
            padding: 2rem;
            border-radius: 10px;
            background-color: #313244;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
            max-width: 800px;
            width: 100%;
        }
        h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            color: #a6e3a1;
        }
        p {
            font-size: 1.2rem;
            margin-bottom: 2rem;
            color: #bac2de;
        }
        .success-icon {
            font-size: 4rem;
            color: #a6e3a1;
            margin-bottom: 1rem;
        }
        .home-button {
            display: inline-block;
            padding: 0.8rem 1.5rem;
            background-color: #74c7ec;
            color: #1e1e2e;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            transition: background-color 0.3s ease;
        }
        .home-button:hover {
            background-color: #89dceb;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✓</div>
        <h1>Subscription Successful!</h1>
        <p>` + username + ` Your subscription will be activated soon.</p>
        <p>You will have full access to our Plex server shortly.</p>
        <p>An invite link will be sent to ` + emailDisplay + `. Please check your inbox (and spam folder).</p>
        <a href="/" class="home-button">Return Home</a>
    </div>
</body>
</html>
`
}

// generateCancelHTML generates the HTML for the cancellation page
func generateCancelHTML(priceID string) string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Subscription Cancelled</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #1e1e2e;
            color: #cdd6f4;
            margin: 0;
            padding: 0;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
        }
        .container {
            text-align: center;
            padding: 2rem;
            border-radius: 10px;
            background-color: #313244;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
            max-width: 800px;
            width: 100%;
        }
        h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            color: #f38ba8;
        }
        p {
            font-size: 1.2rem;
            margin-bottom: 2rem;
            color: #bac2de;
        }
        .cancel-icon {
            font-size: 4rem;
            color: #f38ba8;
            margin-bottom: 1rem;
        }
        .buttons {
            margin-top: 1.5rem;
        }
        .button {
            display: inline-block;
            padding: 0.8rem 1.5rem;
            margin: 0 0.5rem;
            color: #1e1e2e;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            transition: background-color 0.3s ease;
        }
        .home-button {
            background-color: #74c7ec;
        }
        .retry-button {
            background-color: #f9e2af;
        }
        .home-button:hover {
            background-color: #89dceb;
        }
        .retry-button:hover {
            background-color: #f5e0bc;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="cancel-icon">✕</div>
        <h1>Subscription Cancelled</h1>
        <p>Your subscription process was cancelled.</p>
        <p>If you encountered an issue or have changed your mind, you can try again or contact support.</p>
        <div class="buttons">
            <a href="/" class="button home-button">Return Home</a>
            <a href="/api/v1/stripe/checkout?price_id=` + priceID + `" class="button retry-button">Try Again</a>
        </div>
    </div>
</body>
</html>
`
}
