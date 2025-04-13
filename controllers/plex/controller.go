package plexcontroller

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"striplex/config"
	"striplex/model"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PlexController struct {
	basePath string
	client   *http.Client
}

// NewPlexController creates a new PlexController instance.
func NewPlexController(basePath string, client *http.Client) *PlexController {
	return &PlexController{
		basePath: basePath,
		client:   client,
	}
}

func (c *PlexController) GetRoutes(r *gin.RouterGroup) {
	r.GET("/auth", c.Authenticate)
	r.GET("/callback", c.Callback)
}

// PlexPinResponse represents the response from Plex PIN creation API.
type PlexPinResponse struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	AuthToken string `json:"authToken"`
	ClientID  string `json:"clientIdentifier"`
}

// PlexCookieData stores authentication data in the session.
type PlexCookieData struct {
	State string `json:"state"`
	PinID int    `json:"pin_id"`
}

// PlexUserResponse represents the user data returned by Plex API.
type PlexUserResponse struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Title     string `json:"title"`
	Thumb     string `json:"thumb"`
	AuthToken string `json:"authToken"`
}

// Authenticate redirects the user to Plex for authentication.
func (p *PlexController) Authenticate(c *gin.Context) {
	code, err := p.generatePlexPin()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate Plex PIN"})
		return
	}

	// Generate a random state parameter to prevent CSRF attacks
	state := generateRandomState()
	cookieData := PlexCookieData{
		State: state,
		PinID: code.ID,
	}

	jsonData, err := json.Marshal(cookieData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal state data"})
		return
	}

	// Store the state in session/cookie to verify it on callback
	session := sessions.Default(c)
	session.Set("plex_auth_state", string(jsonData))

	if err := session.Save(); err != nil {
		slog.Error("Failed to save session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Build the Plex authentication URL
	baseURL := "https://app.plex.tv/auth#"
	params := url.Values{}
	params.Add("clientID", config.Config.GetString("plex.client_id"))
	params.Add("forwardUrl", fmt.Sprintf("https://%s%s/callback?state=%s&next=%s",
		config.Config.GetString("server.hostname"), p.basePath, state, c.Query("next")))
	params.Add("code", code.Code)
	params.Add("context[device][product]", config.Config.GetString("plex.product"))

	// Redirect user to Plex for authentication
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", baseURL, params.Encode()))
}

// Callback handles the response from Plex after user authentication.
func (p *PlexController) Callback(c *gin.Context) {
	// Get the state from the URL query parameter
	returnedState := c.Query("state")
	session := sessions.Default(c)

	// Get the stored state from cookie
	storedStateStr, exists := session.Get("plex_auth_state").(string)
	if !exists || storedStateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session state"})
		return
	}

	var cookieData PlexCookieData
	if err := json.Unmarshal([]byte(storedStateStr), &cookieData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cookies"})
		return
	}

	// Verify that the state matches to prevent CSRF attacks
	if returnedState != cookieData.State {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	// Check the PIN status to get the auth token
	pinStatus, err := p.checkPlexPin(cookieData.PinID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check PIN status: " + err.Error()})
		return
	}

	// Check if the PIN has been claimed and we have an auth token
	if pinStatus.AuthToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication not completed. Please try again."})
		return
	}

	// Verify the token and get user info
	userInfo, err := p.getPlexUserInfo(pinStatus.AuthToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Failed to verify token: " + err.Error(),
		})
		return
	}

	userInfoData := model.UserInfo{
		ID:       userInfo.ID,
		UUID:     userInfo.UUID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}
	jsonData, err := json.Marshal(userInfoData)
	if err != nil {
		slog.Error("Failed to marshal user info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
		return
	}
	// Clear the state cookie and set user info
	session.Delete("plex_auth_state")
	session.Set("user_info", string(jsonData))

	if err := session.Save(); err != nil {
		slog.Error("Failed to save session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user session"})
		return
	}

	// Check for redirect URL in session (for checkout flow)
	redirectURL := "/"
	// Alternatively check for next parameter in URL
	if nextURL := c.Query("next"); nextURL != "" {
		redirectURL = nextURL
	}

	// Redirect to the appropriate URL
	c.Redirect(http.StatusFound, redirectURL)
}

// generatePlexPin creates a new PIN for PIN-based authentication with Plex.
func (p *PlexController) generatePlexPin() (*PlexPinResponse, error) {
	// Create form data matching the Plex API requirements
	formData := url.Values{}
	formData.Set("strong", "true")
	formData.Set("X-Plex-Product", config.Config.GetString("plex.product"))
	formData.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))

	// Create the request with form data
	req, err := http.NewRequest(
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var pinResponse PlexPinResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinResponse); err != nil {
		return nil, fmt.Errorf("failed to parse PIN response: %w", err)
	}

	return &pinResponse, nil
}

// checkPlexPin checks if the PIN has been claimed and returns the auth token.
func (p *PlexController) checkPlexPin(pinID int) (*PlexPinResponse, error) {
	// Create the request URL with the PIN ID
	reqURL := fmt.Sprintf("https://plex.tv/api/v2/pins/%d", pinID)

	// Create form data
	formData := url.Values{}
	formData.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))

	// Create the request
	req, err := http.NewRequest(
		http.MethodGet,
		reqURL,
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN check request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PIN check request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var pinResponse PlexPinResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinResponse); err != nil {
		return nil, fmt.Errorf("failed to parse PIN check response: %w", err)
	}

	return &pinResponse, nil
}

// getPlexUserInfo verifies the access token and returns user information.
func (p *PlexController) getPlexUserInfo(token string) (*PlexUserResponse, error) {
	// Create the request
	req, err := http.NewRequest(http.MethodGet, "https://plex.tv/api/v2/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Product", config.Config.GetString("plex.product"))
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))
	req.Header.Set("X-Plex-Token", token)

	// Execute the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user info request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token or error getting user info: status code %d", resp.StatusCode)
	}

	// Parse the response
	var userResponse PlexUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	return &userResponse, nil
}

// generateRandomState creates a random string for the state parameter.
func generateRandomState() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// If random generation fails, use a timestamp-based fallback
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
