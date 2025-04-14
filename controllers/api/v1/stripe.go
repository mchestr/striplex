package v1controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"striplex/config"
	"striplex/services"

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
				// Entitlement is being added - generate Wizarr invite code
				wizarrService := services.NewWizarrService(http.DefaultClient)
				invite, err := wizarrService.GenerateInviteLink()
				if err != nil {
					slog.Error("Failed to generate Wizarr invite", "error", err, "customer", summary.Customer)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invite"})
					return
				}

				slog.Info("Generated Wizarr invite code", "code", invite.ID, "customer", summary.Customer)

				// Update customer metadata with the Wizarr invite code
				params := &stripe.CustomerParams{
					Metadata: map[string]string{},
				}
				if existingInvites, ok := stripeCustomer.Metadata["wizarr_invites"]; ok {
					params.Metadata["wizarr_invites"] = existingInvites + "," + strconv.Itoa(invite.ID)
				} else {
					params.Metadata["wizarr_invites"] = strconv.Itoa(invite.ID)
				}

				_, err = customer.Update(summary.Customer, params)
				if err != nil {
					slog.Error("Failed to update customer metadata", "error", err, "customer", summary.Customer)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer metadata"})
					return
				}

				slog.Info("Updated customer with invite code", "customer", summary.Customer, "code", invite.Code, "entitlement", entitlement.LookupKey)
				c.JSON(http.StatusOK, gin.H{"status": "entitlement_added", "type": event.Type, "entitlement": entitlement.LookupKey})
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
		// Entitlement is being removed - log this event
		slog.Info("Entitlement removed", "customer", summary.Customer)

		// Get wizarr_invites from metadata
		if invitesList, ok := stripeCustomer.Metadata["wizarr_invites"]; ok && invitesList != "" {
			inviteIDs := strings.Split(invitesList, ",")
			slog.Info("Found invite IDs to delete", "count", len(inviteIDs), "customer", summary.Customer)

			// Delete each invite
			for _, inviteIDStr := range inviteIDs {
				inviteID, err := strconv.Atoi(strings.TrimSpace(inviteIDStr))
				if err != nil {
					slog.Error("Invalid invite ID format", "id", inviteIDStr, "customer", summary.Customer)
					continue
				}

				// Get Plex username associated with the invitation
				plexUserId, err := s.services.Wizarr.GetPlexIDFromInvitation(inviteID)
				if err != nil {
					slog.Error("Failed to get Plex username from invitation",
						"error", err, "id", inviteID, "customer", summary.Customer)
					continue
				}
				slog.Info("Found Plex username for invitation", "id", inviteID, "username", plexUserId, "customer", summary.Customer)

				if err := s.services.Plex.UnshareLibrary(plexUserId); err != nil {
					slog.Error("Failed to unshare library with Plex user",
						"error", err, "username", plexUserId, "customer", summary.Customer)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unshare library"})
					continue
				}
				slog.Info("Successfully unshared library with Plex user", "id", inviteID, "username", plexUserId, "customer", summary.Customer)

				if err := s.services.Wizarr.DeleteInvite(inviteID); err != nil {
					slog.Error("Failed to delete Wizarr invite",
						"error", err, "id", inviteID, "customer", summary.Customer)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invite"})
					continue
				}
				slog.Info("Successfully deleted Wizarr invite", "id", inviteID, "username", plexUserId, "customer", summary.Customer)

			}

			// Clear the wizarr_invites metadata
			params := &stripe.CustomerParams{
				Metadata: map[string]string{
					"wizarr_invites": "",
				},
			}

			if _, err := customer.Update(summary.Customer, params); err != nil {
				slog.Error("Failed to update customer metadata after deletion",
					"error", err, "customer", summary.Customer)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer metadata"})
				return
			}
		} else {
			slog.Info("No Wizarr invites found for customer", "customer", summary.Customer)
		}
		c.JSON(http.StatusOK, gin.H{"status": "entitlement_removed", "type": event.Type})
	}
}
