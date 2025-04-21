package v1controller

import (
	"net/http"
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
