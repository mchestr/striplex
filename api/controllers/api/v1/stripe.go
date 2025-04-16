package v1controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"plefi/api/config"
	"plefi/api/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"
)

// Webhook handles incoming webhook events from Stripe.
func (s *V1) Webhook(ctx *gin.Context) {
	// Read the request body
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		slog.Error("Failed to read webhook request body", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Get the signature from the header
	sigHeader := ctx.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		slog.Warn("Missing Stripe signature in webhook request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
		return
	}

	// Verify webhook signature and construct the event
	event, err := webhook.ConstructEvent(payload, sigHeader, config.Config.GetString("stripe.webhook_secret"))
	if err != nil {
		slog.Error("Failed to verify webhook signature", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	// Process the webhook event based on its type
	if err := s.processWebhookEvent(ctx, event); err != nil {
		slog.Error("Failed to process webhook event",
			"error", err,
			"event_type", event.Type,
			"event_id", event.ID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetSubscriptions retrieves all subscriptions for the authenticated user
func (s *V1) GetSubscriptions(ctx *gin.Context) {
	// Check user authentication
	session := sessions.Default(ctx)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Authentication required",
		})
		return
	}

	// Parse user info
	var userInfoData models.UserInfo
	if userInfoString, ok := userInfo.(string); ok {
		if err := json.Unmarshal([]byte(userInfoString), &userInfoData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "Invalid session data",
			})
			return
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Unexpected session data format",
		})
		return
	}

	customersIter := customer.Search(
		&stripe.CustomerSearchParams{
			SearchParams: stripe.SearchParams{
				Query: fmt.Sprintf("metadata['plex_user_id']:'%s'", strconv.Itoa(userInfoData.ID)),
			},
		},
	)
	subscriptions := make([]stripe.Subscription, 0)
	for customersIter.Next() {
		cus := customersIter.Customer()

		// Find all subscriptions for this customer
		subscriptionIter := subscription.List(&stripe.SubscriptionListParams{
			Customer: stripe.String(cus.ID),
			Status:   stripe.String("active"),
		})
		for subscriptionIter.Next() {
			sub := subscriptionIter.Subscription()

			// Only include non-empty subscriptions
			if len(sub.Items.Data) == 0 {
				continue
			}
			// Add subscription to results
			subscriptions = append(subscriptions, *sub)
		}

		if err := subscriptionIter.Err(); err != nil {
			slog.Error("Error listing subscriptions", "error", err, "customer", cus.ID)
		}
	}

	if err := customersIter.Err(); err != nil {
		slog.Error("Error listing customers", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to retrieve customer data",
		})
		return
	}

	// Return subscriptions data
	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"subscriptions": subscriptions,
	})
}

// CancelSubscriptionRequest represents the request body for canceling a subscription
type CancelSubscriptionRequest struct {
	SubscriptionID string `json:"subscription_id" binding:"required"`
}

// CancelUserSubscription cancels a specific subscription for the authenticated user
func (s *V1) CancelUserSubscription(ctx *gin.Context) {
	// Check for authentication
	userInfo, err := models.GetUserInfo(ctx)
	if err != nil || userInfo == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Authentication required",
		})
		return
	}

	// Parse request body to get subscription ID
	var reqBody CancelSubscriptionRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid request: subscription_id is required",
		})
		return
	}

	subscription, err := s.services.Stripe.GetSubscription(ctx, userInfo, reqBody.SubscriptionID)
	if err != nil {
		slog.Error("Failed to retrieve subscription",
			"error", err,
			"subscription_id", reqBody.SubscriptionID,
			"user_id", userInfo.ID)
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  "subscription not found",
		})
		return
	}

	// Cancel the specific subscription
	updatedSub, err := s.services.Stripe.CancelAtEndSubscription(ctx, subscription.ID)
	if err != nil {
		slog.Error("Failed to cancel subscription",
			"error", err,
			"subscription_id", reqBody.SubscriptionID,
			"user_id", userInfo.ID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to cancel subscription",
		})
		return
	}

	slog.Info("Subscription canceled",
		"subscription_id", updatedSub.ID,
		"customer_id", updatedSub.Customer.ID,
		"plex_user_id", userInfo.ID)

	// Return success
	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"subscription": updatedSub,
	})
}

// processWebhookEvent handles different types of Stripe webhook events
func (s *V1) processWebhookEvent(ctx context.Context, event stripe.Event) error {
	// Only process entitlements.active_entitlement_summary.updated events
	if event.Type != "entitlements.active_entitlement_summary.updated" {
		slog.Info("Ignoring non-entitlements webhook event", "type", event.Type)
		return nil
	}

	// Parse the event data
	summary, prevAttrs, err := parseEntitlementEventData(event)
	if err != nil {
		return fmt.Errorf("failed to parse webhook event data: %w", err)
	}

	slog.Info("Processing entitlements update",
		"customer_id", summary.Customer,
		"current_count", len(summary.Entitlements.Data),
		"previous_count", len(prevAttrs.Entitlements.Data))

	// Get the Stripe customer details
	stripeCustomer, err := customer.Get(summary.Customer, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve Stripe customer %s: %w", summary.Customer, err)
	}

	// Handle entitlement addition
	if len(summary.Entitlements.Data) > 0 {
		return s.handleEntitlementAddition(ctx, stripeCustomer, summary)
	}

	// Handle entitlement removal
	if len(prevAttrs.Entitlements.Data) > 0 {
		return s.handleEntitlementRemoval(ctx, stripeCustomer)
	}

	// If we reach here, it's an entitlement update that doesn't change the count
	slog.Info("Entitlement updated without count change", "customer", summary.Customer)
	return nil
}

// parseEntitlementEventData extracts the entitlement summary and previous attributes from an event
func parseEntitlementEventData(event stripe.Event) (*stripe.EntitlementsActiveEntitlementSummary, *stripe.EntitlementsActiveEntitlementSummary, error) {
	// Convert event.Data.Object to JSON and then to EntitlementsActiveEntitlementSummary
	rawJSON, err := json.Marshal(event.Data.Object)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	var summary stripe.EntitlementsActiveEntitlementSummary
	if err := json.Unmarshal(rawJSON, &summary); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal event data to summary: %w", err)
	}

	var prevAttrs stripe.EntitlementsActiveEntitlementSummary
	// Process previous attributes if they exist
	if event.Data.PreviousAttributes != nil {
		prevAttrsJSON, err := json.Marshal(event.Data.PreviousAttributes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal previous attributes: %w", err)
		}

		if err := json.Unmarshal(prevAttrsJSON, &prevAttrs); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal previous attributes: %w", err)
		}
	}

	return &summary, &prevAttrs, nil
}

// handleEntitlementAddition processes new entitlements being added to a customer
func (s *V1) handleEntitlementAddition(
	c context.Context,
	stripeCustomer *stripe.Customer,
	summary *stripe.EntitlementsActiveEntitlementSummary,
) error {
	// Iterate through entitlements and handle based on lookup key
	for _, entitlement := range summary.Entitlements.Data {
		switch entitlement.LookupKey {
		case config.Config.GetString("stripe.entitlement_name"):
			return s.shareLibraryWithCustomer(c, stripeCustomer, entitlement)
		default:
			slog.Info("Ignoring entitlement with unsupported lookup key",
				"lookup_key", entitlement.LookupKey,
				"customer", stripeCustomer.ID)
		}
	}

	// If we get here, no matching entitlements were found
	slog.Info("No matching entitlements found", "customer", summary.Customer)
	return nil
}

// shareLibraryWithCustomer shares the Plex library with the customer based on metadata
func (s *V1) shareLibraryWithCustomer(
	ctx context.Context,
	stripeCustomer *stripe.Customer,
	entitlement *stripe.EntitlementsActiveEntitlement,
) error {
	// Get Plex user email from customer metadata
	plexUserEmail, ok := stripeCustomer.Metadata["plex_email"]
	if !ok || plexUserEmail == "" {
		// Fallback to customer email if metadata email is not available
		plexUserEmail = stripeCustomer.Email
		if plexUserEmail == "" {
			return fmt.Errorf("no plex email found for customer %s", stripeCustomer.ID)
		}
		slog.Info("Using customer email instead of metadata", "email", plexUserEmail, "customer", stripeCustomer.ID)
	}

	// Share Plex library with the user
	if err := s.services.Plex.ShareLibrary(ctx, plexUserEmail); err != nil {
		return fmt.Errorf("failed to share Plex library with %s: %w", plexUserEmail, err)
	}

	slog.Info("Shared Plex library with user",
		"customer", stripeCustomer.ID,
		"plex_user", plexUserEmail,
		"entitlement", entitlement.LookupKey)
	return nil
}

// handleEntitlementRemoval processes entitlements being removed from a customer
func (s *V1) handleEntitlementRemoval(
	ctx context.Context,
	stripeCustomer *stripe.Customer,
) error {
	slog.Info("Entitlement removed", "customer", stripeCustomer.ID)

	// First try to get Plex user ID from customer metadata
	plexUserID, ok := stripeCustomer.Metadata["plex_user_id"]
	if !ok || plexUserID == "" {
		slog.Error("No Plex user ID found in customer metadata",
			"customer", stripeCustomer.ID)
		return fmt.Errorf("no plex user ID found for customer %s", stripeCustomer.ID)
	}

	// Unshare library with the Plex user using ID
	if err := s.services.Plex.UnshareLibrary(ctx, plexUserID); err != nil {
		return fmt.Errorf("failed to unshare Plex library with user ID %s: %w", plexUserID, err)
	}

	slog.Info("Successfully unshared library with Plex user", "user_id", plexUserID, "customer", stripeCustomer.ID)
	return nil
}

// parseUserInfo parses user info from the session
func parseUserInfo(userInfo interface{}) (*models.UserInfo, error) {
	var userInfoData models.UserInfo

	if byteData, ok := userInfo.(string); ok {
		if err := json.Unmarshal([]byte(byteData), &userInfoData); err != nil {
			return nil, fmt.Errorf("invalid user info JSON: %w", err)
		}
		return &userInfoData, nil
	}

	return nil, fmt.Errorf("user info is not in expected string format")
}
