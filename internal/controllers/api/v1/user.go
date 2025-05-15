package v1controller

import (
	"log/slog"
	"net/http"
	"plefi/internal/db"
	"plefi/internal/models"

	"github.com/labstack/echo/v4"
)

// GetCurrentUser returns the currently authenticated user's information
func (h *V1) GetCurrentUser(c echo.Context, user *models.UserInfo) error {
	// Return user info
	c.JSON(http.StatusOK, models.GetCurrentUserResponse{
		BaseResponse: models.BaseResponse{
			Status: "success",
		},
		User: user,
	})
	return nil
}

// SetUserNotes sets notes for a user (admin only)
func (h *V1) SetUserNotes(c echo.Context) error {
	// Parse request
	var req models.SetUserNotesRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  "error",
			Message: "Invalid request",
		})
	}

	// Update notes in database
	if err := db.DB.UpdateUserNotes(c.Request().Context(), req.UserID, req.Notes); err != nil {
		slog.Error("Failed to update user notes", "error", err, "user_id", req.UserID)
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Status:  "error",
			Message: "Failed to update user notes",
		})
	}

	// Return success
	return c.JSON(http.StatusOK, models.SetUserNotesResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "User notes updated successfully",
		},
	})
}
