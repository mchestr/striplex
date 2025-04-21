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
