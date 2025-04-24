package v1controller

import (
	"database/sql"
	"log/slog"
	"net/http"
	"plefi/internal/config"
	"plefi/internal/db"
	"plefi/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (h *V1) CheckServerAccess(c echo.Context, user *models.UserInfo) error {
	// Check if the user has access to the server
	hasAccess, err := h.services.Plex.UserHasServerAccess(c.Request().Context(), user.ID)
	if err != nil {
		slog.Error("Failed to check server access",
			"error", err,
			"user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check server access")
	}

	c.JSON(http.StatusOK, models.CheckServerAccessResponse{
		BaseResponse: models.BaseResponse{
			Status: "success",
		},
		HasAccess: hasAccess || user.ID == config.C.Plex.AdminUserID,
	})
	return nil
}

// GetPlexUsersResponse represents the response for listing Plex users
type GetPlexUsersResponse struct {
	models.BaseResponse
	Users []models.PlexUser `json:"users"`
}

// GetPlexUserResponse represents the response for getting a single Plex user
type GetPlexUserResponse struct {
	models.BaseResponse
	User models.PlexUser `json:"user"`
}

// GetPlexUserInvitesResponse represents the response for getting a user's invites
type GetPlexUserInvitesResponse struct {
	models.BaseResponse
	Invites []models.PlexUserInvite `json:"invites"`
}

// RevokeAccessResponse represents the response for revoking a user's access
type RevokeAccessResponse struct {
	models.BaseResponse
}

// GetPlexUsers returns a list of all Plex users (admin only)
func (h *V1) GetPlexUsers(c echo.Context) error {
	// Get all users from the database
	users, err := db.DB.GetAllPlexUsers(c.Request().Context())
	if err != nil {
		slog.Error("Failed to get Plex users", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	return c.JSON(http.StatusOK, GetPlexUsersResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Users retrieved successfully",
		},
		Users: users,
	})
}

// GetPlexUser returns details of a specific Plex user (admin only)
func (h *V1) GetPlexUser(c echo.Context) error {
	// Get user ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Get user from the database
	user, err := db.DB.GetPlexUser(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		slog.Error("Failed to get Plex user", "error", err, "user_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user")
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, GetPlexUserResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "User retrieved successfully",
		},
		User: *user,
	})
}

// GetPlexUserInvites returns all invites for a specific Plex user (admin only)
func (h *V1) GetPlexUserInvites(c echo.Context) error {
	// Get user ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Check if user exists
	user, err := db.DB.GetPlexUser(c.Request().Context(), id)
	if err != nil || user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Get user's invites from the database
	invites, err := db.DB.GetPlexUserInvites(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to get user invites", "error", err, "user_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user invites")
	}

	return c.JSON(http.StatusOK, GetPlexUserInvitesResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Invites retrieved successfully",
		},
		Invites: invites,
	})
}

// RevokePlexAccess revokes a user's access to the Plex server (admin only)
func (h *V1) RevokePlexAccess(c echo.Context) error {
	// Get user ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Check if user exists
	user, err := db.DB.GetPlexUser(c.Request().Context(), id)
	if err != nil || user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Don't allow revoking admin access
	if user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot revoke access for admin users")
	}

	// Revoke access in Plex
	if err := h.services.Plex.UnshareLibrary(c.Request().Context(), id); err != nil {
		slog.Error("Failed to revoke Plex access", "error", err, "user_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke Plex access")
	}

	return c.JSON(http.StatusOK, RevokeAccessResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "access revoked successfully",
		},
	})
}

// DeleteCurrentUser deletes the currently authenticated user's information
func (h *V1) DeletePlexUser(c echo.Context, user *models.UserInfo) error {
	if config.C.Plex.AdminUserID == user.ID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete admin users")
	} // Get user ID from path parameter

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}
	// Unshare the Plex library with the user
	if err := h.services.Plex.UnshareLibrary(c.Request().Context(), id); err != nil {
		slog.Error("Failed to unshare Plex library with user",
			"error", err,
			"user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unshare Plex library")
	}

	// Delete the user from the database
	if err := db.DB.DeletePlexUser(c.Request().Context(), id); err != nil {
		slog.Error("Failed to delete user from database",
			"error", err,
			"user_id", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user")
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Status:  "success",
		Message: "user deleted successfully",
	})
}
