package plexcontroller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"plefi/internal/config"
	"plefi/internal/db"
	"plefi/internal/models"
	"plefi/internal/services"
	"plefi/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
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
func (c *PlexController) GetRoutes(r *echo.Group) {
	r.GET("/auth", c.Authenticate)
	r.GET("/callback", c.Callback)
}

type PlexAuthenticateRequest struct {
	Next string `query:"next"`
}

// Authenticate redirects the user to Plex for authentication.
func (h *PlexController) Authenticate(c echo.Context) error {
	var req PlexAuthenticateRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		return err
	}
	code, err := h.services.Plex.GeneratePin(c.Request().Context())
	if err != nil {
		slog.Error("Failed to generate Plex PIN", "error", err)
		return err
	}

	authData := &models.PlexAuth{
		State: generateRandomState(),
		PinID: code.ID,
	}
	if err := utils.SaveSessionData(c, utils.PlexSessionState, authData); err != nil {
		slog.Error("Failed to save Plex auth data to session", "error", err)
		return err
	}
	c.Redirect(http.StatusTemporaryRedirect, buildPlexAuthURL(h.basePath, authData.State, req.Next, code.Code))
	return nil
}

type PlexCallbackRequest struct {
	Next  string `query:"next"`
	State string `query:"state"`
}

// Callback handles the response from Plex after user authentication.
func (h *PlexController) Callback(c echo.Context) error {
	var req PlexCallbackRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind callback request", "error", err)
		return err
	}

	plexAuth, err := utils.GetSessionData(c, utils.PlexSessionState)
	if err != nil {
		slog.Error("Session state not found")
		return fmt.Errorf("session state not found")
	}
	plexAuthData, ok := plexAuth.(*models.PlexAuth)
	if !ok {
		slog.Error("Invalid session state type")
		return fmt.Errorf("invalid session state type")
	}
	if plexAuthData.State != req.State {
		slog.Error("State mismatch", "expected", plexAuthData.State, "received", req.State)
		return fmt.Errorf("state mismatch")
	}

	// Check the PIN status to get the auth token
	pinStatus, err := h.services.Plex.ClaimPin(c.Request().Context(), plexAuthData.PinID)
	if err != nil {
		slog.Error("Failed to check PIN status", "error", err, "pin_id", plexAuthData.PinID)
		return err
	}

	// Check if the PIN has been claimed and we have an auth token
	if pinStatus.AuthToken == "" {
		slog.Warn("Authentication not completed", "pin_id", plexAuthData.PinID)
		return fmt.Errorf("PIN not claimed yet")
	}

	// Verify the token and get user info
	userInfo, err := h.services.Plex.GetUserDetails(c.Request().Context(), pinStatus.AuthToken)
	if err != nil {
		slog.Error("Failed to verify token", "error", err)
		return err
	}
	if err := utils.SaveSessionData(c, utils.UserInfoState, &models.UserInfo{
		ID:       userInfo.ID,
		UUID:     userInfo.UUID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		IsAdmin:  config.C.Plex.AdminUserID == userInfo.ID,
	}); err != nil {
		slog.Error("Failed to save user info to session", "error", err)
		return err
	}
	if err := utils.SaveSessionData(c, utils.PlexSessionState, nil); err != nil {
		slog.Error("Failed to clear plex auth from session", "error", err)
		return err
	}
	if err := db.DB.SavePlexUser(c.Request().Context(), models.PlexUser{
		ID:       userInfo.ID,
		UUID:     userInfo.UUID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		IsAdmin:  config.C.Plex.AdminUserID == userInfo.ID,
	}); err != nil {
		slog.Error("Failed to save Plex user to database", "error", err)
		return err
	}

	if err := db.DB.SavePlexToken(c.Request().Context(), models.PlexToken{
		UserID:      userInfo.ID,
		AccessToken: pinStatus.AuthToken,
	}); err != nil {
		slog.Error("Failed to save Plex token to database", "error", err)
		return err
	}

	// Determine where to redirect after successful authentication
	redirectURL := "/"
	if req.Next != "" {
		redirectURL = req.Next
	}

	// Redirect to the appropriate URL
	c.Redirect(http.StatusFound, redirectURL)
	return nil
}

// buildPlexAuthURL constructs the Plex authentication URL
func buildPlexAuthURL(basePath, state, nextURL, code string) string {
	baseURL := "https://app.plex.tv/auth#"
	params := url.Values{}
	params.Add("clientID", config.C.Plex.ClientID)
	params.Add("forwardUrl", fmt.Sprintf("https://%s%s/callback?state=%s&next=%s",
		config.C.Server.Hostname, basePath, state, nextURL))
	params.Add("code", code)
	params.Add("context[device][product]", config.C.Plex.ProductName)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
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
