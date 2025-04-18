package v1controller

import (
	"log/slog"
	"net/http"
	"plefi/api/config"
	"plefi/api/models"
	"plefi/api/utils"

	"github.com/labstack/echo/v4"
)

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (h *V1) CheckServerAccess(c echo.Context) error {
	// Get user info from the session
	userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		return err
	}
	if userInfo == nil {
		c.JSON(http.StatusOK, map[string]any{
			"status":     "success",
			"has_access": false,
		})
		return nil
	}
	userInfoData, ok := userInfo.(*models.UserInfo)
	if !ok {
		c.JSON(http.StatusOK, map[string]any{
			"status":     "success",
			"has_access": false,
		})
		return nil
	}

	if userInfoData.ID == config.C.Plex.AdminUserID {
		c.JSON(http.StatusOK, map[string]any{
			"status":     "success",
			"has_access": true,
		})
		return nil
	}

	// Check if the user has access to the server
	hasAccess, err := h.services.Plex.UserHasServerAccess(c.Request().Context(), userInfoData.ID)
	if err != nil {
		slog.Error("Failed to check server access",
			"error", err,
			"user_id", userInfoData.ID)
		return err
	}

	c.JSON(http.StatusOK, map[string]any{
		"status":     "success",
		"has_access": hasAccess,
	})
	return nil
}
