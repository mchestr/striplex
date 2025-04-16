package plexcontroller

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"plefi/api/config"
	"plefi/api/models"
	"plefi/api/services"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PlexController handles Plex API authentication and user management
type PlexController struct {
	basePath string
	client   *http.Client
	services *services.Services
}

// NewPlexController creates a new PlexController instance.
func NewPlexController(basePath string, client *http.Client, services *services.Services) *PlexController {
	return &PlexController{
		basePath: basePath,
		client:   client,
		services: services,
	}
}

// GetRoutes registers all routes for the Plex controller
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
func (p *PlexController) Authenticate(ctx *gin.Context) {
	code, err := p.generatePlexPin(ctx)
	if err != nil {
		slog.Error("Failed to generate Plex PIN", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate Plex PIN"})
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
		slog.Error("Failed to marshal state data", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal state data"})
		return
	}

	// Store the state in session/cookie to verify it on callback
	session := sessions.Default(ctx)
	session.Set("plex_auth_state", string(jsonData))

	if err := session.Save(); err != nil {
		slog.Error("Failed to save session", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Get next URL for redirect after authentication
	nextURL := ctx.Query("next")

	// Build the Plex authentication URL
	authURL := buildPlexAuthURL(p.basePath, state, nextURL, code.Code)

	// Redirect user to Plex for authentication
	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles the response from Plex after user authentication.
func (p *PlexController) Callback(ctx *gin.Context) {
	// Get the state from the URL query parameter
	returnedState := ctx.Query("state")
	session := sessions.Default(ctx)

	// Validate stored session state
	cookieData, err := validateSessionState(session, returnedState)
	if err != nil {
		slog.Error("Session state validation failed", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check the PIN status to get the auth token
	pinStatus, err := p.checkPlexPin(ctx, cookieData.PinID)
	if err != nil {
		slog.Error("Failed to check PIN status", "error", err, "pin_id", cookieData.PinID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check PIN status: " + err.Error()})
		return
	}

	// Check if the PIN has been claimed and we have an auth token
	if pinStatus.AuthToken == "" {
		slog.Warn("Authentication not completed", "pin_id", cookieData.PinID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Authentication not completed. Please try again."})
		return
	}

	// Verify the token and get user info
	userInfo, err := p.getPlexUserInfo(ctx, pinStatus.AuthToken)
	if err != nil {
		slog.Error("Failed to verify token", "error", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Failed to verify token: " + err.Error(),
		})
		return
	}

	// Save user information to session
	if err := saveUserSession(session, userInfo); err != nil {
		slog.Error("Failed to save user session", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user session"})
		return
	}

	// Determine where to redirect after successful authentication
	redirectURL := "/"
	if nextURL := ctx.Query("next"); nextURL != "" {
		redirectURL = nextURL
	}

	// Redirect to the appropriate URL
	ctx.Redirect(http.StatusFound, redirectURL)
}

// generatePlexPin creates a new PIN for PIN-based authentication with Plex.
func (p *PlexController) generatePlexPin(ctx context.Context) (*PlexPinResponse, error) {
	// Create form data matching the Plex API requirements
	formData := url.Values{}
	formData.Set("strong", "true")
	formData.Set("X-Plex-Product", config.Config.GetString("plex.product"))
	formData.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))

	// Create the request with form data
	req, err := http.NewRequestWithContext(
		ctx,
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

// checkPlexPin checks if the PIN has been claimed and returns the auth token.
func (p *PlexController) checkPlexPin(ctx context.Context, pinID int) (*PlexPinResponse, error) {
	// Create the request URL with the PIN ID
	reqURL := fmt.Sprintf("https://plex.tv/api/v2/pins/%d", pinID)

	// Create form data
	formData := url.Values{}
	formData.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))

	// Create the request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil, // GET requests don't need a body
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN check request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Client-Identifier", config.Config.GetString("plex.client_id"))

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

// getPlexUserInfo verifies the access token and returns user information.
func (p *PlexController) getPlexUserInfo(ctx context.Context, token string) (*PlexUserResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token provided")
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://plex.tv/api/v2/user", nil)
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid token or error getting user info: status code %d, body: %s",
			resp.StatusCode, string(body))
	}

	// Parse the response
	var userResponse PlexUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	return &userResponse, nil
}

// Helper functions

// buildPlexAuthURL constructs the Plex authentication URL
func buildPlexAuthURL(basePath, state, nextURL, code string) string {
	baseURL := "https://app.plex.tv/auth#"
	params := url.Values{}
	params.Add("clientID", config.Config.GetString("plex.client_id"))
	params.Add("forwardUrl", fmt.Sprintf("https://%s%s/callback?state=%s&next=%s",
		config.Config.GetString("server.hostname"), basePath, state, nextURL))
	params.Add("code", code)
	params.Add("context[device][product]", config.Config.GetString("plex.product"))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// validateSessionState validates the session state against the returned state
func validateSessionState(session sessions.Session, returnedState string) (*PlexCookieData, error) {
	storedStateStr, exists := session.Get("plex_auth_state").(string)
	if !exists || storedStateStr == "" {
		return nil, fmt.Errorf("invalid session state")
	}

	var cookieData PlexCookieData
	if err := json.Unmarshal([]byte(storedStateStr), &cookieData); err != nil {
		return nil, fmt.Errorf("invalid cookies")
	}

	// Verify that the state matches to prevent CSRF attacks
	if returnedState != cookieData.State {
		return nil, fmt.Errorf("invalid state parameter")
	}

	return &cookieData, nil
}

// saveUserSession saves the user information to the session
func saveUserSession(session sessions.Session, userInfo *PlexUserResponse) error {
	userInfoData := models.UserInfo{
		ID:       userInfo.ID,
		UUID:     userInfo.UUID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}
	jsonData, err := json.Marshal(userInfoData)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %w", err)
	}

	// Clear the state cookie and set user info
	session.Delete("plex_auth_state")
	session.Set("user_info", string(jsonData))

	if err := session.Save(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
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
