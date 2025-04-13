package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"striplex/config"
)

// WizarrService handles interactions with the Wizarr API
type WizarrService struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// WizarrInviteRequest represents the request body for creating an invite
type WizarrInviteRequest struct {
	ExpiresAfter int    `json:"expires_after"` // Time in seconds until the invite expires
	MaxUses      int    `json:"max_uses"`      // Maximum number of uses for the invite
	ProfileID    string `json:"profile_id"`    // Optional: The profile ID to assign to the user
}

// WizarrInviteResponse represents the response from the Wizarr API when creating an invite
type WizarrInviteResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Link    string `json:"link"`
	Error   string `json:"error,omitempty"`
}

// NewWizarrService creates a new WizarrService instance
func NewWizarrService(baseURL string, client *http.Client) *WizarrService {
	return &WizarrService{
		client:  client,
		baseURL: baseURL,
		apiKey:  config.Config.GetString("wizarr.api_key"),
	}
}

// GenerateInviteLink creates a new invite in Wizarr and returns the generated link
func (w *WizarrService) GenerateInviteLink(expiresAfter, maxUses int) (string, error) {
	inviteReq := WizarrInviteRequest{
		ExpiresAfter: expiresAfter,
		MaxUses:      maxUses,
	}

	jsonBody, err := json.Marshal(inviteReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal invite request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/invite", w.baseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", w.apiKey)

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var inviteResp WizarrInviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&inviteResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors
	if !inviteResp.Success {
		if inviteResp.Error != "" {
			return "", fmt.Errorf("wizarr API error: %s", inviteResp.Error)
		}
		return "", fmt.Errorf("wizarr API returned unsuccessful response")
	}

	return inviteResp.Link, nil
}
