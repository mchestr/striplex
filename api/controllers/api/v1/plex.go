package v1controller

import (
	"log/slog"
	"net/http"
	"plefi/api/config"
	"plefi/api/models"

	"github.com/gin-gonic/gin"
)

// CheckServerAccess checks if the authenticated user has access to the Plex server
func (c *V1) CheckServerAccess(ctx *gin.Context) {
	// Get user info from the session
	userInfo, err := models.GetUserInfo(ctx)
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Internal server error",
		})
		return
	}

	if userInfo == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"has_access": false,
		})
		return
	}

	if userInfo.ID == config.C.Plex.AdminUserID {
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"has_access": true,
		})
		return
	}

	// Check if the user has access to the server
	hasAccess, err := c.services.Plex.UserHasServerAccess(ctx, userInfo.ID)
	if err != nil {
		slog.Error("Failed to check server access",
			"error", err,
			"user_id", userInfo.ID)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to check server access",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"has_access": hasAccess,
	})
}
