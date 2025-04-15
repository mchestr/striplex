package model

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetUserInfo(ctx *gin.Context) (*UserInfo, error) {
	session := sessions.Default(ctx)
	userInfo := session.Get("user_info")
	if userInfo == nil {
		return nil, nil
	}

	var userInfoData UserInfo
	if byteData, ok := userInfo.(string); ok {
		if err := json.Unmarshal([]byte(byteData), &userInfoData); err != nil {
			return nil, fmt.Errorf("invalid user info JSON: %w", err)
		}
		return &userInfoData, nil
	}
	return nil, fmt.Errorf("user info is not in expected string format")
}
