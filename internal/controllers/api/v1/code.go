package v1controller

import (
	"log/slog"
	"math/rand"
	"net/http"
	"plefi/internal/db"
	"plefi/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// CreateInviteCodeRequest represents the request body for creating an invite code
type CreateInviteCodeRequest struct {
	Code            string     `json:"code"`
	MaxUses         *int       `json:"max_uses"`
	ExpiresAt       *time.Time `json:"expires_at"`
	EntitlementName string     `json:"entitlement_name"`
	Duration        *time.Time `json:"duration"`
}

// CreateInviteCodeResponse represents the response for create invite code request
type CreateInviteCodeResponse struct {
	models.BaseResponse
	InviteCode models.InviteCode `json:"invite_code"`
}

// ListInviteCodesResponse represents the response for list invite codes request
type ListInviteCodesResponse struct {
	models.BaseResponse
	InviteCodes []models.InviteCode `json:"invite_codes"`
}

// DeleteInviteCodeResponse represents the response for delete invite code request
type DeleteInviteCodeResponse struct {
	models.BaseResponse
}

// GetInviteCodeUsersResponse represents the response for getting users of an invite code
type GetInviteCodeResponse struct {
	models.BaseResponse
	models.InviteCode
	Users []models.PlexUser `json:"users"`
}

// ClaimInviteCodeRequest represents the request body for claiming an invite code
type ClaimInviteCodeRequest struct {
	Code string `json:"code"`
}

// ClaimInviteCodeResponse represents the response for claim invite code request
type ClaimInviteCodeResponse struct {
	models.BaseResponse
	InviteCode models.InviteCode `json:"invite_code"`
}

// generateRandomCode creates a random invite code of specified length
func generateRandomCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[r.Intn(len(charset))])
	}

	return sb.String()
}

// CreateInviteCode creates a new invite code
func (h *V1) CreateInviteCode(c echo.Context) error {
	// Parse request body
	var req CreateInviteCodeRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Generate random code if not provided
	if req.Code == "" {
		req.Code = generateRandomCode(8) // Generate an 8-character code
	}

	// Validation
	if req.EntitlementName == "" {
		req.EntitlementName = "plex" // Default entitlement name
	}

	// Create invite code model
	inviteCode := models.InviteCode{
		Code:            req.Code,
		MaxUses:         req.MaxUses,
		ExpiresAt:       req.ExpiresAt,
		EntitlementName: req.EntitlementName,
		Duration:        req.Duration,
		UsedCount:       0,
		IsDisabled:      false,
	}

	// Save to database
	id, err := db.DB.SaveInviteCode(c.Request().Context(), inviteCode)
	if err != nil {
		slog.Error("Failed to save invite code", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create invite code")
	}

	// Set the ID in the response
	inviteCode.ID = id

	// Return success response
	return c.JSON(http.StatusCreated, CreateInviteCodeResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Invite code created successfully",
		},
		InviteCode: inviteCode,
	})
}

// GetInviteCodeUsers retrieves all Plex users who have used a specific invite code
func (h *V1) GetInviteCode(c echo.Context) error {
	// Get code ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invite code ID")
	}

	code, err := db.DB.GetInviteCode(c.Request().Context(), id)
	if err != nil || code == nil {
		slog.Error("Failed to get invite code", "error", err, "code_id", id)
		return echo.NewHTTPError(http.StatusNotFound, "code not found")
	}

	// Get users who have used this invite code from database
	users, err := db.DB.GetUsersWithActiveInviteCode(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to get users for invite code", "error", err, "code_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users for invite code")
	}

	// Return success response
	return c.JSON(http.StatusOK, GetInviteCodeResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "code retrieved successfully",
		},
		InviteCode: *code,
		Users:      users,
	})
}

// ListInviteCodes lists all active invite codes
func (h *V1) ListInviteCodes(c echo.Context) error {
	// Get active invite codes from database
	codes, err := db.DB.ListActiveInviteCodes(c.Request().Context())
	if err != nil {
		slog.Error("Failed to list invite codes", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve invite codes")
	}

	// Return success response
	return c.JSON(http.StatusOK, ListInviteCodesResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Invite codes retrieved successfully",
		},
		InviteCodes: codes,
	})
}

// ClaimInviteCode allows a Plex user to claim an invite code
func (h *V1) ClaimInviteCode(c echo.Context, user *models.UserInfo) error {
	// Parse request body
	var req ClaimInviteCodeRequest
	if err := c.Bind(&req); err != nil || req.Code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Get all active codes and find the matching one
	inviteCode, err := db.DB.GetInviteCodeByCode(c.Request().Context(), req.Code)
	if err != nil || inviteCode == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "code not found")
	}

	now := time.Now()
	if inviteCode.IsDisabled ||
		(inviteCode.Duration != nil && inviteCode.Duration.Before(now)) ||
		(inviteCode.ExpiresAt != nil && inviteCode.ExpiresAt.Before(now)) ||
		(inviteCode.UsedCount >= *inviteCode.MaxUses) {
		slog.Error("invite code is disabled or expired",
			"code_id", inviteCode.ID,
			"code", req.Code,
			"duration", inviteCode.Duration,
			"expires_at", inviteCode.ExpiresAt,
			"used_count", inviteCode.UsedCount,
			"max_uses", *inviteCode.MaxUses)
		return echo.NewHTTPError(http.StatusBadRequest, "code not found")
	}

	// Associate the code with the user
	err = db.DB.AssociatePlexUserWithInviteCode(c.Request().Context(), user.ID, inviteCode.ID)
	if err != nil {
		slog.Error("Failed to associate user with invite code", "error", err, "user_id", user.ID, "code_id", inviteCode.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to claim invite code")
	}

	// Update the usage count of the invite code
	err = db.DB.UpdateInviteCodeUsage(c.Request().Context(), inviteCode.ID)
	if err != nil {
		slog.Error("Failed to update invite code usage", "error", err, "code_id", inviteCode.ID)
		// Not returning an error here as the code was already claimed
	}

	// Get the user's details to get email address
	plexUser, err := db.DB.GetPlexUser(c.Request().Context(), user.ID)
	if err != nil {
		slog.Error("Failed to get user details", "error", err, "user_id", user.ID)
		// Continue despite error, as the code was claimed successfully
	} else if plexUser != nil && plexUser.Email != "" {
		// Share the Plex library with the user
		invite, err := h.services.Plex.ShareLibrary(c.Request().Context(), plexUser.Email)
		if err != nil {
			slog.Error("Failed to share Plex library with user", "error", err, "user_id", user.ID, "email", plexUser.Email)
			// Continue despite error, as the code was claimed successfully
		} else {
			slog.Info("Plex library shared with user", "user_id", user.ID, "email", plexUser.Email, "invite_id", invite.ID)

			// Try to auto-accept the invitation if we have the user's Plex token
			token, tokenErr := db.DB.GetPlexToken(c.Request().Context(), user.ID)
			if tokenErr == nil && token != nil && token.AccessToken != "" {
				if acceptErr := h.services.Plex.AcceptInvite(c.Request().Context(), token.AccessToken, invite.ID); acceptErr != nil {
					slog.Error("Failed to auto-accept Plex invite", "error", acceptErr, "user_id", user.ID, "invite_id", invite.ID)
				} else {
					slog.Info("Plex invite auto-accepted", "user_id", user.ID, "invite_id", invite.ID)
				}
			}
		}
	}

	// Return success response
	return c.JSON(http.StatusOK, ClaimInviteCodeResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Invite code claimed successfully",
		},
		InviteCode: *inviteCode,
	})
}

// DeleteInviteCode disables an existing invite code
func (h *V1) DeleteInviteCode(c echo.Context) error {
	// Get code ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invite code ID")
	}

	// Disable the invite code in database
	err = db.DB.DisableInviteCode(c.Request().Context(), id)
	if err != nil {
		slog.Error("Failed to disable invite code", "error", err, "code_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to disable invite code")
	}

	// Return success response
	return c.JSON(http.StatusOK, DeleteInviteCodeResponse{
		BaseResponse: models.BaseResponse{
			Status:  "success",
			Message: "Invite code disabled successfully",
		},
	})
}
