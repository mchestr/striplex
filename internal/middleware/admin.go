package middleware

import (
	"log/slog"
	"net/http"
	"plefi/internal/models"
	"plefi/internal/utils"

	"github.com/labstack/echo/v4"
)

func NewAdminMiddleware(adminID int) echo.MiddlewareFunc {
	// AdminMiddleware ensures the current user is an admin.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user info from the session
			userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
			if err != nil || userInfo == nil {
				slog.Error("Failed to get user info", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			userInfoData, ok := userInfo.(*models.UserInfo)
			if !ok || userInfoData.ID == adminID {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			return next(c)
		}
	}
}
