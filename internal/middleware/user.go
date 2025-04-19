package middleware

import (
	"log/slog"
	"net/http"
	"plefi/internal/models"
	"plefi/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserHandlerFunc func(c echo.Context, user *models.UserInfo) error

func UserHandler(handler UserHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user info from the session
		userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
		if err != nil || userInfo == nil {
			slog.Error("Failed to get user info", "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		userInfoData, ok := userInfo.(*models.UserInfo)
		if !ok {
			slog.Error("Failed to cast user info to UserInfo type")
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return handler(c, userInfoData)
	}
}

func AnonymousHandler(handler UserHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user info from the session
		userInfo, err := utils.GetSessionData(c, utils.UserInfoState)
		if err != nil || userInfo == nil {
			slog.Error("Failed to get user info", "error", err)
			return handler(c, nil)
		}
		userInfoData, ok := userInfo.(*models.UserInfo)
		if !ok {
			slog.Error("Failed to cast user info to UserInfo type")
			return handler(c, nil)
		}
		return handler(c, userInfoData)
	}
}
