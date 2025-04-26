package v1controller

import (
	"database/sql"
	"log/slog"
	"net/http"
	"plefi/internal/config"
	"plefi/internal/db"
	"plefi/internal/models"
	"plefi/internal/services/plex"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (h *V1) GetServerAccess(c echo.Context, user *models.UserInfo) error {
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

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (h *V1) CheckServerAccess(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}
	// Check if the user has access to the server
	hasAccess, err := h.services.Plex.UserHasServerAccess(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to check server access",
			"error", err,
			"user_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check server access")
	}

	c.JSON(http.StatusOK, models.CheckServerAccessResponse{
		BaseResponse: models.BaseResponse{
			Status: "success",
		},
		HasAccess: hasAccess || id == config.C.Plex.AdminUserID,
	})
	return nil
}

// GetPlexUsersResponse represents the response for listing Plex users
type GetPlexUsersResponse struct {
	models.BaseResponse
	Users []models.PlexUserWithAccess `json:"users"`
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

// GrantAccessResponse represents the response for granting a user's access
type GrantAccessResponse struct {
	models.BaseResponse
}

// GrantPlexAccessRequest represents the request body for granting Plex access
type GrantPlexAccessRequest struct {
	UserID string `json:"user_id"`
}

// GetPlexUsers returns a list of all Plex users (admin only)
func (h *V1) GetPlexUsers(c echo.Context) error {
	// Get all users from the database
	users, err := db.DB.GetAllPlexUsers(c.Request().Context())
	if err != nil {
		slog.Error("Failed to get Plex users", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	plexUsers, err := h.services.Plex.GetUsers(c.Request().Context())
	if err != nil {
		slog.Error("Failed to get Plex users from API", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users from Plex API")
	}
	//Create map for fast lookup
	plexUserMap := make(map[int]plex.PlexUser)
	for _, plexUser := range plexUsers {
		plexUserMap[plexUser.ID] = plexUser
	}

	plexUsersWithAccess := make([]models.PlexUserWithAccess, len(users))
	for i, user := range users {
		plexUsersWithAccess[i] = models.PlexUserWithAccess{
			PlexUser:  user,
			HasAccess: h.services.Plex.CheckUserHasAccess(plexUserMap, user.ID),
		}
	}

	return c.JSON(http.StatusOK, GetPlexUsersResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Users retrieved successfully",
		},
		Users: plexUsersWithAccess,
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

// GrantPlexAccess grants a user access to the Plex server (admin only)
func (h *V1) GrantPlexAccess(c echo.Context) error {
	// Parse request body
	var req GrantPlexAccessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	id, err := strconv.Atoi(req.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Check if user exists
	user, err := db.DB.GetPlexUser(c.Request().Context(), id)
	if err != nil || user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Check if user already has access
	hasAccess, err := h.services.Plex.UserHasServerAccess(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to check server access", "error", err, "user_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check server access")
	}

	if hasAccess {
		return echo.NewHTTPError(http.StatusBadRequest, "User already has access")
	}

	// User needs an email to grant access
	if user.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User has no email address")
	}

	// Share Plex library with the user
	invite, err := h.services.Plex.ShareLibrary(c.Request().Context(), user.Email)
	if err != nil {
		slog.Error("Failed to share Plex library with user", "error", err, "user_id", id, "email", user.Email)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to grant Plex access")
	}

	// Try to auto-accept the invitation if we have the user's Plex token
	token, tokenErr := db.DB.GetPlexToken(c.Request().Context(), id)
	if tokenErr == nil && token != nil && token.AccessToken != "" {
		if acceptErr := h.services.Plex.AcceptInvite(c.Request().Context(), token.AccessToken, invite.ID); acceptErr != nil {
			slog.Error("Failed to auto-accept invite", "error", acceptErr, "invite_id", invite.ID)
			// Continue despite error as the invite was sent successfully
		} else {
			slog.Info("Auto-accepted invite for user", "user_id", id, "invite_id", invite.ID)
		}
	}

	slog.Info("Plex library shared with user", "user_id", id, "email", user.Email, "invite_id", invite.ID)
	return c.JSON(http.StatusOK, GrantAccessResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "access granted successfully",
		},
	})
}

// DeleteCurrentUser deletes the currently authenticated user's information
func (h *V1) DeletePlexUser(c echo.Context, user *models.UserInfo) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}
	if config.C.Plex.AdminUserID == id {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete admin users")
	}

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
