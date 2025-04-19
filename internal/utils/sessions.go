package utils

import (
	"plefi/internal/config"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSessionData(c echo.Context, key string) (interface{}, error) {
	s, err := session.Get(config.C.Auth.SessionName, c)
	if err != nil {
		return nil, err
	}
	return s.Values[key], nil
}

func SaveSessionData(c echo.Context, key string, data interface{}) error {
	s, err := session.Get(config.C.Auth.SessionName, c)
	if err != nil {
		return err
	}
	s.Values[key] = data
	return s.Save(c.Request(), c.Response())
}
