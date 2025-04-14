package services

import (
	"fmt"
	"net/http"
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
