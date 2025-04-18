package v1controller

import (
	"net/http"
	"plefi/api/models"
	"plefi/api/utils"

	"log/slog"

	"github.com/labstack/echo/v4"
)

// GetCurrentUser returns the currently authenticated user's information
func (h *V1) GetCurrentUser(c echo.Context) error {
	userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		return err
	}
	if userInfo == nil {
		// User is not authenticated
		c.JSON(http.StatusOK, map[string]any{
			"status": "success",
			"user":   nil,
		})
		return nil
	}

	userInfoData, ok := userInfo.(*models.UserInfo)
	if !ok {
		c.JSON(http.StatusOK, map[string]any{
			"status": "success",
			"user":   nil,
		})
		return nil
	}

	// Return user info
	c.JSON(http.StatusOK, map[string]any{
		"status": "success",
		"user": map[string]any{
			"id":       userInfoData.ID,
			"uuid":     userInfoData.UUID,
			"username": userInfoData.Username,
			"email":    userInfoData.Email,
		},
	})
	return nil
}
