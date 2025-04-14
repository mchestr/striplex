package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"plefi/config"
	"strings"
)

// PlexService handles interactions with the Plex Media Server API
type PlexService struct {
	client   *http.Client
	token    string
	serverID string
}

// NewPlexService creates a new PlexService instance
func NewPlexService(client *http.Client) *PlexService {
	return &PlexService{
		client:   client,
		token:    config.Config.GetString("plex.token"),
		serverID: config.Config.GetString("plex.server_id"),
	}
}

// PlexUser represents a user in the Plex system
type PlexUser struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Title      string `json:"title"`
	Thumb      string `json:"thumb"`
	HomeSize   int    `json:"homeSize"`
	AllowSync  bool   `json:"allowSync"`
	AllowTuner int    `json:"allowTuners"`
}

// PlexFriendsResponse represents the response structure when fetching friends list
type PlexFriendsResponse struct {
	MediaContainer struct {
		Size      int        `json:"size"`
		Users     []PlexUser `json:"User"`
		PublicKey string     `json:"publicKey"`
	} `json:"MediaContainer"`
}

// PlexLibrarySection represents a library section in a Plex server
type PlexLibrarySection struct {
	ID    int    `json:"id"`
	Key   int    `json:"key"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// PlexServerResponse represents the JSON response from the Plex server API
type PlexServerResponse struct {
	Name            string               `json:"name"`
	MachineID       string               `json:"machineIdentifier"`
	LibrarySections []PlexLibrarySection `json:"librarySections"`
}

// plexError represents a structured error response from Plex API
type plexError struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	} `json:"errors"`
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
func (p *PlexService) ShareLibrary(ctx context.Context, email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	sectionIDs, err := p.GetSectionIDsByNames(ctx, strings.Split(config.Config.GetString("plex.shared_libraries"), ","))
	if err != nil {
		return fmt.Errorf("failed to get section IDs: %w", err)
	}

	payload := map[string]interface{}{
		"invitedEmail":      email,
		"machineIdentifier": p.serverID,
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
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://clients.plex.tv/api/v2/shared_servers",
		bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create share request: %w", err)
	}

	p.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("share request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status code
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusUnauthorized:
		slog.Warn("Unauthorized Plex token", "email", email)
		return fmt.Errorf("invalid Plex token")
	case http.StatusBadRequest:
		// Try to parse structured error if available
		var plexErr plexError
		if err := json.Unmarshal(body, &plexErr); err == nil && len(plexErr.Errors) > 0 {
			return fmt.Errorf("bad request: %s", plexErr.Errors[0].Message)
		}
		return fmt.Errorf("bad request: %s", string(body))
	default:
		slog.Debug("Share library failed",
			"status", resp.Status,
			"response", string(body),
			"email", email)
		return fmt.Errorf("API returned unexpected status: %d %s", resp.StatusCode, resp.Status)
	}
}

// GetSectionIDsByNames retrieves section IDs that match the provided section names
func (p *PlexService) GetSectionIDsByNames(ctx context.Context, sectionNames []string) ([]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://plex.tv/api/v2/servers/%s", p.serverID), nil)
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

// setCommonHeaders sets the common headers used in Plex API requests
func (p *PlexService) setCommonHeaders(req *http.Request) {
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))
	req.Header.Set("Accept", "application/json")
}
