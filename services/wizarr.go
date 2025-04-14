package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	Duration          int      `json:"duration,omitempty"`           // Time in seconds until the invite expires
	MaxUses           int      `json:"max_uses,omitempty"`           // Maximum number of uses for the invite
	Expires           int      `json:"expires"`                      // Time in seconds until the invite expires
	Unlimited         bool     `json:"unlimited"`                    // Whether the invite has unlimited uses
	PlexAllowSync     bool     `json:"plex_allow_sync"`              // Whether Plex sync is allowed
	LiveTV            bool     `json:"live_tv"`                      // Whether Live TV is allowed
	HideUser          bool     `json:"hide_user"`                    // Whether to hide the user
	AllowDownload     bool     `json:"allow_download"`               // Whether downloads are allowed
	Sessions          int      `json:"sessions"`                     // Number of sessions allowed
	SpecificLibraries []string `json:"specific_libraries,omitempty"` // List of library IDs
}

// WizarrInviteResponse represents the response from the Wizarr API when creating an invite
type WizarrInviteResponse struct {
	AllowDownload     bool      `json:"allow_download"`
	Code              string    `json:"code"`
	Created           string    `json:"created"`
	Duration          string    `json:"duration"` // Changed from *int to string
	Expires           string    `json:"expires"`  // Changed from *int to string
	HideUser          bool      `json:"hide_user"`
	ID                int       `json:"id"`
	LiveTV            *bool     `json:"live_tv"` // Changed to pointer to handle null
	PlexAllowSync     bool      `json:"plex_allow_sync"`
	Sessions          *int      `json:"sessions"`
	SpecificLibraries *[]string `json:"specific_libraries"` // Changed to pointer to handle null
	Unlimited         bool      `json:"unlimited"`
	Used              bool      `json:"used"`
	UsedAt            *string   `json:"used_at"`
	UsedBy            []string  `json:"used_by"` // Changed to []string from *string
}

// WizarrUser represents a user in the Wizarr system
type WizarrUser struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Auth     *string `json:"auth"`
	Code     *string `json:"code"`
	Token    string  `json:"token"`
	Expires  *string `json:"expires"`
	Created  string  `json:"created"`
}

// NewWizarrService creates a new WizarrService instance
func NewWizarrService(client *http.Client) *WizarrService {
	return &WizarrService{
		client:  client,
		baseURL: config.Config.GetString("wizarr.url"),
		apiKey:  config.Config.GetString("wizarr.api_key"),
	}
}

// ListInvites retrieves all invitations from the Wizarr API
func (w *WizarrService) ListInvites() ([]WizarrInviteResponse, error) {
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/invitations", w.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse response
	var invites []WizarrInviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&invites); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return invites, nil
}

// ListUsers retrieves all users from the Wizarr API
func (w *WizarrService) ListUsers() ([]WizarrUser, error) {
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users", w.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse response
	var users []WizarrUser
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return users, nil
}

// GetUserByEmail finds a user by their email address
func (w *WizarrService) GetUserByEmail(email string) (*WizarrUser, error) {
	// Get all users first
	users, err := w.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Find the user with the matching email
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

// GetUserByEmail finds a user by their email address
func (w *WizarrService) GetUserId(id int) (*WizarrUser, error) {
	// Get all users first
	users, err := w.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Find the user with the matching email
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with id %d not found", id)
}

// DeleteInvite deletes an invitation by its ID
func (w *WizarrService) DeleteInvite(id int) error {
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/invitations/%d", w.baseURL, id), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
}

// GenerateInviteLink creates a new invite in Wizarr and returns the generated link
func (w *WizarrService) GenerateInviteLink() (*WizarrInviteResponse, error) {
	// Create invite request using the struct
	inviteReq := WizarrInviteRequest{
		Expires: int(config.Config.GetDuration("wizarr.invite_expiration").Minutes()),
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(inviteReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/invitations", w.baseURL),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var inviteResp WizarrInviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&inviteResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &inviteResp, nil
}

// GetInviteByID retrieves a specific invitation by its ID
func (w *WizarrService) GetInviteByID(id int) (*WizarrInviteResponse, error) {
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/invitations/%d", w.baseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))

	// Send the request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse response
	var invite WizarrInviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &invite, nil
}

// GetUsernameFromInvitation retrieves the username associated with an invitation
func (w *WizarrService) GetPlexIDFromInvitation(inviteId int) (string, error) {
	invite, err := w.GetInviteByID(inviteId)
	if err != nil {
		return "", fmt.Errorf("failed to get invite by ID: %w", err)
	}
	// Check if the invite has been used
	if !invite.Used || len(invite.UsedBy) == 0 {
		return "", fmt.Errorf("invitation has not been used yet")
	}

	// Get the email from the UsedBy field (first entry)
	userIdentifier := invite.UsedBy[0]

	// Try by id first
	if userId, err := strconv.Atoi(userIdentifier); err == nil {
		user, err := w.GetUserId(userId)
		if err != nil {
			return "", fmt.Errorf("failed to find user by ID or email: %w", err)
		}
		return user.Token, nil
	}
	// Find the user by email
	user, err := w.GetUserByEmail(userIdentifier)
	if err != nil {
	}
	return user.Token, nil
}
