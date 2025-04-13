package controllers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"striplex/config"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82/webhook"
)

// StripeController handles Stripe webhook events.
type StripeController struct {
	basePath string
}

// NewStripeController creates a new StripeController instance.
func NewStripeController(basePath string) *StripeController {
	return &StripeController{
		basePath: basePath,
	}
}

// Entitlements represents the structure of entitlement data from Stripe.
type Entitlements struct {
	Customer     string `json:"customer"`
	Entitlements struct {
		Data []struct {
			LookupKey string `json:"lookup_key"`
		} `json:"data"`
	} `json:"entitlements"`
}

// Webhook handles incoming webhook events from Stripe.
func (s *StripeController) Webhook(c *gin.Context) {
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

	// Handle events related to entitlements
	var entitlement Entitlements

	// Convert the raw JSON to our struct
	jsonBytes, err := json.Marshal(event.Data.Object)
	if err != nil {
		slog.Error("Error marshaling event data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing event data"})
		return
	}

	if err := json.Unmarshal(jsonBytes, &entitlement); err != nil {
		slog.Error("Error unmarshaling entitlements data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing entitlements data"})
		return
	}

	slog.Info("Processed entitlements", "customer", entitlement.Customer, "entitlements_count", len(entitlement.Entitlements.Data))

	// Return a success response to Stripe
	c.JSON(http.StatusOK, gin.H{"status": "success", "type": event.Type})
}
