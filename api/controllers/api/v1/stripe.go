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
	"plefi/api/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"
)

// Webhook handles incoming webhook events from Stripe.
func (h *V1) Webhook(c echo.Context) error {
	// Read the request body
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		slog.Error("Failed to read webhook request body", "error", err)
		return err
	}

	// Get the signature from the header
	sigHeader := c.Request().Header.Get("Stripe-Signature")
	if sigHeader == "" {
		slog.Warn("Missing Stripe signature in webhook request")
		return fmt.Errorf("missing Stripe signature")
	}

	// Verify webhook signature and construct the event
	event, err := webhook.ConstructEvent(payload, sigHeader, config.C.Stripe.WebhookSecret.Value())
	if err != nil {
		slog.Error("Failed to verify webhook signature", "error", err)
		return err
	}

	// Process the webhook event based on its type
	if err := h.processWebhookEvent(c.Request().Context(), event); err != nil {
		slog.Error("Failed to process webhook event",
			"error", err,
			"event_type", event.Type,
			"event_id", event.ID)
		return err
	}
	c.JSON(http.StatusOK, map[string]any{"status": "success"})
	return nil
}

// SubscriptionItem represents a simplified subscription item
type SubscriptionItem struct {
	ID               string `json:"id"`
	CurrentPeriodEnd int64  `json:"current_period_end"`
	Price            struct {
		UnitAmount int64  `json:"unit_amount"`
		Currency   string `json:"currency"`
		Recurring  struct {
			Interval string `json:"interval"`
		} `json:"recurring"`
		Product struct {
			Name string `json:"name"`
		} `json:"product"`
	} `json:"price"`
}

// SimplifiedSubscription represents minimal subscription data needed by frontend
type SimplifiedSubscription struct {
	ID                string             `json:"id"`
	Status            string             `json:"status"`
	CancelAtPeriodEnd bool               `json:"cancel_at_period_end"`
	Items             []SubscriptionItem `json:"items"`
}

// GetSubscriptions retrieves all subscriptions for the authenticated user
func (h *V1) GetSubscriptions(c echo.Context) error {
	// Check user authentication
	userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
	if err != nil || userInfo == nil {
		slog.Error("Failed to get user info", "error", err)
		return fmt.Errorf("user not authenticated")
	}
	userInfoData, ok := userInfo.(*models.UserInfo)
	if !ok {
		slog.Error("Failed to cast user info to UserInfo type")
		return fmt.Errorf("failed to cast user info to UserInfo type")
	}

	customersIter := customer.Search(
		&stripe.CustomerSearchParams{
			SearchParams: stripe.SearchParams{
				Query: fmt.Sprintf("metadata['plex_user_id']:'%s'", strconv.Itoa(userInfoData.ID)),
			},
		},
	)

	simplifiedSubs := make([]SimplifiedSubscription, 0)

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

			// Create simplified subscription
			simplifiedSub := SimplifiedSubscription{
				ID:                sub.ID,
				Status:            string(sub.Status),
				CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
				Items:             make([]SubscriptionItem, 0, len(sub.Items.Data)),
			}

			// Add simplified items
			for _, item := range sub.Items.Data {
				simpleItem := SubscriptionItem{
					ID:               item.ID,
					CurrentPeriodEnd: item.CurrentPeriodEnd,
				}

				if item.Price != nil {
					simpleItem.Price.UnitAmount = item.Price.UnitAmount
					simpleItem.Price.Currency = string(item.Price.Currency)

					if item.Price.Recurring != nil {
						simpleItem.Price.Recurring.Interval = string(item.Price.Recurring.Interval)
					}

					if item.Price.Product != nil {
						simpleItem.Price.Product.Name = item.Price.Product.Name
					}
				}

				simplifiedSub.Items = append(simplifiedSub.Items, simpleItem)
			}

			simplifiedSubs = append(simplifiedSubs, simplifiedSub)
		}

		if err := subscriptionIter.Err(); err != nil {
			slog.Error("Error listing subscriptions", "error", err, "customer", cus.ID)
		}
	}

	if err := customersIter.Err(); err != nil {
		return err
	}

	// Return subscriptions data
	c.JSON(http.StatusOK, map[string]any{
		"status":        "success",
		"subscriptions": simplifiedSubs,
	})
	return nil
}

// CancelSubscriptionRequest represents the request body for canceling a subscription
type CancelSubscriptionRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

// CancelUserSubscription cancels a specific subscription for the authenticated user
func (h *V1) CancelUserSubscription(c echo.Context) error {
	// Check for authentication
	userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
	if err != nil {
		return err
	}
	if userInfo == nil {
		return fmt.Errorf("user not authenticated")
	}
	userInfoData, ok := userInfo.(*models.UserInfo)
	if !ok {
		return fmt.Errorf("failed to cast user info to UserInfo type")
	}

	// Parse request body to get subscription ID
	var reqBody CancelSubscriptionRequest
	if err := c.Bind(&reqBody); err != nil {
		return err
	}

	subscription, err := h.services.Stripe.GetSubscription(c.Request().Context(), userInfoData, reqBody.SubscriptionID)
	if err != nil {
		slog.Error("Failed to retrieve subscription",
			"error", err,
			"subscription_id", reqBody.SubscriptionID,
			"user_id", userInfoData.ID)
		return err
	}

	// Cancel the specific subscription
	updatedSub, err := h.services.Stripe.CancelAtEndSubscription(c.Request().Context(), subscription.ID)
	if err != nil {
		slog.Error("Failed to cancel subscription",
			"error", err,
			"subscription_id", reqBody.SubscriptionID,
			"user_id", userInfoData.ID)
		return err
	}

	slog.Info("Subscription canceled",
		"subscription_id", updatedSub.ID,
		"customer_id", updatedSub.Customer.ID,
		"plex_user_id", userInfoData.ID)

	// Return success
	c.JSON(http.StatusOK, map[string]any{
		"status":       "success",
		"subscription": updatedSub,
	})
	return nil
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
		case config.C.Stripe.EntitlementName:
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
