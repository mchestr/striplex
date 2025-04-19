package v1controller

import (
	"log/slog"
	"net/http"
	"plefi/internal/config"
	"plefi/internal/models"

	"github.com/labstack/echo/v4"
)

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (h *V1) CheckServerAccess(c echo.Context, user *models.UserInfo) error {
	if user.ID == config.C.Plex.AdminUserID {
		c.JSON(http.StatusOK, map[string]any{
			"status":     "success",
			"has_access": true,
		})
		return nil
	}

	// Check if the user has access to the server
	hasAccess, err := h.services.Plex.UserHasServerAccess(c.Request().Context(), user.ID)
	if err != nil {
		slog.Error("Failed to check server access",
			"error", err,
			"user_id", user.ID)
		return err
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":     "success",
		"has_access": hasAccess,
	})
	return nil
}
