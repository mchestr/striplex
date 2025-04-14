package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"striplex/config"
)

// PlexService handles interactions with the Plex Media Server API
type PlexService struct {
	client   *http.Client
	baseURL  string
	token    string
	serverId string
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

// PlexServerResponse represents the JSON response from the Plex server API
type PlexServerResponse struct {
	Name            string `json:"name"`
	MachineID       string `json:"machineIdentifier"`
	LibrarySections []struct {
		ID    int    `json:"id"`
		Key   int    `json:"key"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"librarySections"`
}

// NewPlexService creates a new PlexService instance
func NewPlexService(client *http.Client) *PlexService {
	return &PlexService{
		client:   client,
		baseURL:  config.Config.GetString("plex.url"),
		token:    config.Config.GetString("plex.token"),
		serverId: config.Config.GetString("plex.server_id"),
	}
}

// UnshareLibrary removes a user's access to the Plex server
func (p *PlexService) UnshareLibrary(userID string) error {
	url := fmt.Sprintf("https://plex.tv/api/v2/sharings/%s", userID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// ShareLibrary shares specific libraries with a Plex user
func (p *PlexService) ShareLibrary(email string) error {
	// Create the payload with the new structure
	sectionIDs, err := p.GetSectionIDsByNames(strings.Split(config.Config.GetString("plex.shared_libraries"), ","))
	if err != nil {
		return fmt.Errorf("failed to get section IDs: %w", err)
	}
	payload := map[string]interface{}{
		"invitedEmail":      email,
		"machineIdentifier": p.serverId,
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

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, "https://clients.plex.tv/api/v2/shared_servers", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid Plex token")
	case http.StatusBadRequest:
		// Read the response body for more detailed error information
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("bad request and failed to read error details: %w", err)
		}
		return fmt.Errorf("bad request: %s", string(body))
	default:
		return fmt.Errorf("API returned unexpected status: %d %s", resp.StatusCode, resp.Status)
	}
}

// GetSectionIDsByNames retrieves section IDs that match the provided section names
func (p *PlexService) GetSectionIDsByNames(sectionNames []string) ([]int, error) {
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://plex.tv/api/v2/servers/%s", p.serverId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse the JSON response
	var serverInfo PlexServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverInfo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Find matching section IDs
	sectionNameToId := make(map[string]int)
	for _, section := range serverInfo.LibrarySections {
		sectionNameToId[section.Title] = section.ID
	}
	sectionIDs := make([]int, 0, len(sectionNames))
	for _, name := range sectionNames {
		if id, ok := sectionNameToId[name]; ok {
			sectionIDs = append(sectionIDs, id)
		}
	}

	if len(sectionIDs) == 0 {
		return nil, fmt.Errorf("no matching library sections found")
	}

	return sectionIDs, nil
}
