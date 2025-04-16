package v1controller

import (
	"net/http"
	"plefi/api/models"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser returns the currently authenticated user's information
func (c *V1) GetCurrentUser(ctx *gin.Context) {
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
		// User is not authenticated
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   nil,
		})
		return
	}

	// Return user info
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"id":       userInfo.ID,
			"uuid":     userInfo.UUID,
			"username": userInfo.Username,
			"email":    userInfo.Email,
		},
	})
}
