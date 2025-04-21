package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"plefi/internal/config"
	"strings"
)

// PlexServicer defines the interface for Plex Media Server API operations
type PlexServicer interface {
	// UnshareLibrary removes a user's access to the Plex server
	UnshareLibrary(ctx context.Context, userID string) error

	// ShareLibrary shares specific libraries with a Plex user
	ShareLibrary(ctx context.Context, email string) (*PlexShareResponse, error)

	// GetSectionIDsByNames retrieves section IDs that match the provided section names
	GetSectionIDsByNames(ctx context.Context, sectionNames []string) ([]int, error)

	// GetUsers retrieves all users associated with the Plex server
	GetUsers(ctx context.Context) ([]PlexUser, error)

	// UserHasServerAccess checks if a user is in the users list and has access to the server
	UserHasServerAccess(ctx context.Context, userID int) (bool, error)

	// GetUserDetails retrieves detailed information about the authenticated user
	GetUserDetails(ctx context.Context, plexToken string) (*PlexDetailedUserResponse, error)

	// ClaimPin retrieves the PIN status from Plex API
	ClaimPin(ctx context.Context, pinID int) (*PlexPinResponse, error)

	// GeneratePin creates a new Plex PIN for user authentication
	GeneratePin(ctx context.Context) (*PlexPinResponse, error)

	// AcceptInvite allows a user to accept a Plex library invitation using the invite token
	AcceptInvite(ctx context.Context, plexToken string, inviteID int) error

	// GetMachineIdentity returns the server's machineIdentifier from the identity endpoint
	GetMachineIdentity(ctx context.Context, url, plexToken string) (string, error)
}

// Verify that PlexService implements the PlexServicer interface
var _ PlexServicer = (*PlexService)(nil)

// PlexService handles interactions with the Plex Media Server API
type PlexService struct {
	client *http.Client
	token  string
}

// NewPlexService creates a new PlexService instance
func NewPlexService(client *http.Client) *PlexService {
	return &PlexService{
		client: client,
		token:  config.C.Plex.Token.Value(),
	}
}

// UnshareLibrary removes a user's access to the Plex server
func (p *PlexService) UnshareLibrary(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	url := fmt.Sprintf("https://plex.tv/api/v2/sharings/%s", userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create unshare request: %w", err)
	}

	p.setCommonHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("unshare request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		slog.Debug("Unshare library failed",
			"status", resp.Status,
			"response", string(body),
			"userID", userID)
		return fmt.Errorf("unshare API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// ShareLibrary shares specific libraries with a Plex user
func (p *PlexService) ShareLibrary(ctx context.Context, email string) (*PlexShareResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	sectionIDs, err := p.GetSectionIDsByNames(ctx, config.C.Plex.SharedLibraries)
	if err != nil {
		return nil, fmt.Errorf("failed to get section IDs: %w", err)
	}

	payload := map[string]interface{}{
		"invitedEmail":      email,
		"machineIdentifier": config.C.Plex.MachineIdentifier,
		"librarySectionIds": sectionIDs,
		"skipFriendship":    true,
		"settings": map[string]interface{}{
			"allowSync":          false,
			"allowChannels":      false,
			"allowSubtitleAdmin": false,
			"allowTuners":        0,
			"filterMovies":       "",
			"filterMusic":        "",
			"filterPhotos":       "",
			"filterTelevision":   "",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://clients.plex.tv/api/v2/shared_servers",
		bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create share request: %w", err)
	}

	p.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("share request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status code
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		var shareResp PlexShareResponse
		if err := json.Unmarshal(body, &shareResp); err != nil {
			return nil, fmt.Errorf("failed to parse share response: %w", err)
		}
		return &shareResp, nil
	case http.StatusUnauthorized:
		slog.Warn("Unauthorized Plex token", "email", email)
		return nil, fmt.Errorf("invalid Plex token")
	case http.StatusBadRequest:
		// Try to parse structured error if available
		var plexErr PlexErrorResponse
		if err := json.Unmarshal(body, &plexErr); err == nil && len(plexErr.Errors) > 0 {
			return nil, fmt.Errorf("bad request: %s", plexErr.Errors[0].Message)
		}
		return nil, fmt.Errorf("bad request: %s", string(body))
	default:
		slog.Debug("Share library failed",
			"status", resp.Status,
			"response", string(body),
			"email", email)
		return nil, fmt.Errorf("API returned unexpected status: %d %s", resp.StatusCode, resp.Status)
	}
}

// GetUsers retrieves all users associated with the Plex server
func (p *PlexService) GetUsers(ctx context.Context) ([]PlexUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://clients.plex.tv/api/users", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.setCommonHeaders(req)
	// Override the Accept header to ensure we get XML response
	req.Header.Set("Accept", "application/xml")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to get users failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Debug("Get users failed",
			"status", resp.Status,
			"response", string(body))
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	var usersResponse PlexUsersResponse
	if err := xml.Unmarshal(body, &usersResponse); err != nil {
		slog.Debug("Failed to unmarshal XML response",
			"error", err,
			"response_sample", string(body[:min(500, len(body))]))
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return usersResponse.Users, nil
}

// GetSectionIDsByNames retrieves section IDs that match the provided section names
func (p *PlexService) GetSectionIDsByNames(ctx context.Context, sectionNames []string) ([]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://plex.tv/api/v2/servers/%s", config.C.Plex.MachineIdentifier), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	p.setCommonHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Debug("Get section IDs failed",
			"status", resp.Status,
			"response", string(body))
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	var serverInfo PlexServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverInfo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Find matching section IDs
	sectionNameToID := make(map[string]int, len(serverInfo.LibrarySections))
	for _, section := range serverInfo.LibrarySections {
		sectionNameToID[strings.ToLower(section.Title)] = section.ID
	}

	sectionIDs := make([]int, 0, len(sectionNames))
	for _, name := range sectionNames {
		if id, ok := sectionNameToID[strings.ToLower(name)]; ok {
			sectionIDs = append(sectionIDs, id)
		}
	}
	return sectionIDs, nil
}

// UserHasServerAccess checks if a user is in the users list and has access to the server
func (p *PlexService) UserHasServerAccess(ctx context.Context, userID int) (bool, error) {
	users, err := p.GetUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get users: %w", err)
	}

	serverID := config.C.Plex.MachineIdentifier
	for _, user := range users {
		if user.ID == userID {
			// Found the user, now check if they have access to our server
			for _, server := range user.Servers {
				if server.MachineIdentifier == serverID {
					return true, nil
				}
			}
			// User found but doesn't have access to our server
			return false, nil
		}
	}

	// User not found
	return false, nil
}

// GetUserDetails retrieves detailed information about the authenticated user
func (p *PlexService) GetUserDetails(ctx context.Context, plexToken string) (*PlexDetailedUserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://plex.tv/api/v2/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user details request: %w", err)
	}
	p.setCommonHeaders(req)
	req.Header.Set("X-Plex-Token", plexToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user details request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Debug("Get user details failed",
			"status", resp.Status,
			"response", string(body))
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	var userDetails PlexDetailedUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userDetails); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &userDetails, nil
}

func (p *PlexService) ClaimPin(c context.Context, pinID int) (*PlexPinResponse, error) {
	// Create the request URL with the PIN ID
	reqURL := fmt.Sprintf("https://plex.tv/api/v2/pins/%d", pinID)

	// Create form data
	formData := url.Values{}
	formData.Set("X-Plex-Client-Identifier", config.C.Plex.ClientID)

	// Create the request with context
	req, err := http.NewRequestWithContext(
		c,
		http.MethodGet,
		reqURL,
		nil, // GET requests don't need a body
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN check request: %w", err)
	}
	p.setCommonHeaders(req)

	// Execute the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PIN check request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var pinResponse PlexPinResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinResponse); err != nil {
		return nil, fmt.Errorf("failed to parse PIN check response: %w", err)
	}

	return &pinResponse, nil
}

func (p *PlexService) GeneratePin(c context.Context) (*PlexPinResponse, error) {
	// Create form data matching the Plex API requirements
	formData := url.Values{
		"strong":                   {"true"},
		"X-Plex-Product":           {config.C.Plex.ProductName},
		"X-Plex-Client-Identifier": {config.C.Plex.ClientID},
	}

	// Create the request with form data
	req, err := http.NewRequestWithContext(
		c,
		http.MethodPost,
		"https://plex.tv/api/v2/pins",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PIN request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var pinResponse PlexPinResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinResponse); err != nil {
		return nil, fmt.Errorf("failed to parse PIN response: %w", err)
	}

	return &pinResponse, nil
}

// AcceptInvite allows a user to accept a Plex library invitation using the invite token
func (p *PlexService) AcceptInvite(ctx context.Context, plexToken string, inviteID int) error {
	reqURL := fmt.Sprintf("https://plex.tv/api/v2/shared_servers/%d/accept", inviteID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create accept invite request: %w", err)
	}

	p.setCommonHeaders(req)
	req.Header.Set("X-Plex-Token", plexToken)
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("accept invite request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		slog.Debug("Accept invite failed",
			"status", resp.Status,
			"response", string(body))
		return fmt.Errorf("accept invite API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// GetMachineIdentity retrieves the machineIdentifier from the Plex identity endpoint
func (p *PlexService) GetMachineIdentity(ctx context.Context, url, plexToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/identity", url), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create identity request: %w", err)
	}
	p.setCommonHeaders(req)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", plexToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("identity request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status: %d %s", resp.StatusCode, string(body))
	}

	var ir struct {
		MediaContainer struct {
			MachineIdentifier string `json:"machineIdentifier"`
		} `json:"MediaContainer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return "", fmt.Errorf("failed to parse identity response: %w", err)
	}

	return ir.MediaContainer.MachineIdentifier, nil
}

// setCommonHeaders sets the common headers used in Plex API requests
func (p *PlexService) setCommonHeaders(req *http.Request) {
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("X-Plex-Client-Identifier", config.C.Plex.ClientID)
	req.Header.Set("Accept", "application/json")
}

// Helper function to return the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
