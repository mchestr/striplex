package v1controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"striplex/config"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/webhook"
)

// Webhook handles incoming webhook events from Stripe.
func (s *V1) Webhook(c *gin.Context) {
	// Read the request body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.Error("Error reading request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Get the signature from the header
	sigHeader := c.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
		return
	}

	// Verify and construct the event
	event, err := webhook.ConstructEvent(payload, sigHeader, config.Config.GetString("stripe.webhook_secret"))
	if err != nil {
		slog.Error("Error verifying webhook signature", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	// Only process entitlements.created or entitlements.updated events
	if event.Type != "entitlements.active_entitlement_summary.updated" {
		slog.Info("Ignoring non-entitlements webhook event", "type", event.Type)
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "type": event.Type})
		return
	}

	// Handle events related to entitlements
	// Convert event.Data.Object to JSON and then to EntitlementsActiveEntitlementSummary
	rawJSON, err := json.Marshal(event.Data.Object)
	if err != nil {
		slog.Error("Failed to marshal event data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data format"})
		return
	}

	var summary stripe.EntitlementsActiveEntitlementSummary
	if err := json.Unmarshal(rawJSON, &summary); err != nil {
		slog.Error("Failed to unmarshal event data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data format"})
		return
	}

	var prevAttrs struct {
		Entitlements stripe.EntitlementsActiveEntitlementList `json:"entitlements"`
	}
	// Process previous attributes if they exist
	if event.Data.PreviousAttributes != nil {
		prevAttrsJSON, err := json.Marshal(event.Data.PreviousAttributes)
		if err != nil {
			slog.Error("Failed to marshal previous attributes", "error", err)
		} else {
			if err := json.Unmarshal(prevAttrsJSON, &prevAttrs); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data format"})
				return
			}
		}
	}

	slog.Info("Processing entitlements",
		"customer", summary.Customer,
		"current_count", len(summary.Entitlements.Data),
		"previous_count", len(prevAttrs.Entitlements.Data))

	stripeCustomer, err := customer.Get(summary.Customer, nil)
	if err != nil {
		slog.Error("Failed to retrieve Stripe customer", "error", err, "customer", summary.Customer)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer"})
		return
	}

	// Check if an entitlement is being added
	if len(summary.Entitlements.Data) > 0 {
		// Iterate through entitlements and handle based on lookup key
		for _, entitlement := range summary.Entitlements.Data {
			switch entitlement.LookupKey {
			case config.Config.GetString("stripe.entitlement_name"):
				// Get Plex user ID from customer metadata or email
				plexUserEmail, ok := stripeCustomer.Metadata["plex_email"]
				if !ok {
					slog.Error("No Plex user ID found for customer", "customer", summary.Customer)
					c.JSON(http.StatusBadRequest, gin.H{"error": "No Plex user ID found"})
					return
				}

				// Share Plex library with the user
				if err := s.services.Plex.ShareLibrary(plexUserEmail); err != nil {
					slog.Error("Failed to share Plex library", "error", err, "user", plexUserEmail, "customer", summary.Customer)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share Plex library"})
					return
				}

				slog.Info("Shared Plex library with user", "customer", summary.Customer, "plex_user", plexUserEmail, "entitlement", entitlement.LookupKey)
				c.JSON(http.StatusOK, gin.H{"status": "library_shared", "type": event.Type, "entitlement": entitlement.LookupKey})
				return
			default:
				slog.Info("Ignoring entitlement with unsupported lookup key", "lookup_key", entitlement.LookupKey)
			}
		}

		// If we get here, no matching entitlements were found
		slog.Info("No matching entitlements found", "customer", summary.Customer)
		c.JSON(http.StatusOK, gin.H{"status": "no_matching_entitlements", "type": event.Type})
		return
	}

	// Check if an entitlement is being removed
	if len(prevAttrs.Entitlements.Data) > 0 {
		// Entitlement is being removed - unshare Plex library
		slog.Info("Entitlement removed", "customer", summary.Customer)

		// Get Plex user ID from customer metadata
		plexUserID, ok := stripeCustomer.Metadata["plex_user_id"]
		if !ok {
			slog.Error("No Plex user ID found for customer", "customer", summary.Customer)
			c.JSON(http.StatusOK, gin.H{"status": "No Plex user ID found"})
			return
		}

		// Unshare library with the Plex user
		if err := s.services.Plex.UnshareLibrary(plexUserID); err != nil {
			slog.Error("Failed to unshare library with Plex user",
				"error", err, "user", plexUserID, "customer", summary.Customer)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unshare Plex library"})
			return
		}

		slog.Info("Successfully unshared library with Plex user", "user", plexUserID, "customer", summary.Customer)
		c.JSON(http.StatusOK, gin.H{"status": "library_unshared", "type": event.Type})
		return
	}

	// If we reach here, it's an entitlement update that doesn't change the count
	slog.Info("Entitlement updated without count change", "customer", summary.Customer)
	c.JSON(http.StatusOK, gin.H{"status": "entitlement_updated", "type": event.Type})
}
